package app

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xiaojiu/cliplink/internal/static"
)

type staticSite struct {
	webFS fs.FS
	pages map[string]string
}

// SetupStaticRoutes registers the embedded Next.js static export.
func SetupStaticRoutes(router *gin.Engine) {
	site := newStaticSite(static.GetWebFS())
	site.Register(router)
	router.NoRoute(site.Handle)
}

// StaticFileHandler is kept for focused tests and custom routers.
func StaticFileHandler(webFS fs.FS) gin.HandlerFunc {
	return newStaticSite(webFS).Handle
}

func newStaticSite(webFS fs.FS) *staticSite {
	return &staticSite{
		webFS: webFS,
		pages: discoverExportedPages(webFS),
	}
}

func (s *staticSite) Register(router *gin.Engine) {
	routes := make([]string, 0, len(s.pages))
	for routePath := range s.pages {
		routes = append(routes, routePath)
	}
	sort.Strings(routes)

	for _, routePath := range routes {
		routePath := routePath
		handler := func(c *gin.Context) {
			if s.serveRoute(c, c.Request.URL.Path) {
				return
			}
			s.serveNotFound(c)
		}

		router.GET(routePath, handler)
		router.HEAD(routePath, handler)
	}
}

func (s *staticSite) Handle(c *gin.Context) {
	if isAPIPath(c.Request.URL.Path) {
		c.Next()
		return
	}

	if s.serveRoute(c, c.Request.URL.Path) {
		return
	}

	s.serveNotFound(c)
}

func (s *staticSite) serveRoute(c *gin.Context, requestPath string) bool {
	filePath, ok := normalizeStaticPath(requestPath)
	if !ok {
		return false
	}

	for _, candidate := range s.candidateFiles(filePath) {
		if serveFile(c, s.webFS, candidate) {
			return true
		}
	}

	return false
}

func (s *staticSite) candidateFiles(filePath string) []string {
	if filePath == "" {
		filePath = "index.html"
	}

	if routeFile, ok := s.pages["/"+strings.Trim(filePath, "/")]; ok {
		return []string{routeFile}
	}

	candidates := []string{filePath}
	if path.Ext(filePath) == "" {
		candidates = append(candidates, filePath+".html", path.Join(filePath, "index.html"))
	}

	return candidates
}

func (s *staticSite) serveNotFound(c *gin.Context) {
	if serveFileWithStatus(c, s.webFS, "404.html", http.StatusNotFound) {
		return
	}
	c.String(http.StatusNotFound, http.StatusText(http.StatusNotFound))
}

func discoverExportedPages(webFS fs.FS) map[string]string {
	pages := make(map[string]string)

	_ = fs.WalkDir(webFS, ".", func(filePath string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() || path.Ext(filePath) != ".html" {
			return nil
		}

		routePath, ok := exportedPageRoute(filePath)
		if ok {
			pages[routePath] = filePath
		}

		return nil
	})

	return pages
}

func exportedPageRoute(filePath string) (string, bool) {
	if filePath == "404.html" {
		return "", false
	}

	if path.Base(filePath) == "index.html" {
		dir := path.Dir(filePath)
		if dir == "." {
			return "/", true
		}
		return "/" + dir, true
	}

	return "/" + strings.TrimSuffix(filePath, ".html"), true
}

func isAPIPath(requestPath string) bool {
	return requestPath == "/api" || strings.HasPrefix(requestPath, "/api/")
}

func normalizeStaticPath(requestPath string) (string, bool) {
	if requestPath == "" || requestPath == "/" {
		return "index.html", true
	}

	normalized := strings.ReplaceAll(requestPath, "\\", "/")
	for _, segment := range strings.Split(normalized, "/") {
		if segment == ".." {
			return "", false
		}
	}

	filePath := strings.TrimPrefix(path.Clean("/"+strings.TrimPrefix(normalized, "/")), "/")
	if filePath == "" || filePath == "." {
		return "index.html", true
	}

	return filePath, fs.ValidPath(filePath)
}

func serveFile(c *gin.Context, webFS fs.FS, filePath string) bool {
	return serveFileWithStatus(c, webFS, filePath, http.StatusOK)
}

func serveFileWithStatus(c *gin.Context, webFS fs.FS, filePath string, statusCode int) bool {
	file, err := webFS.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil || stat.IsDir() {
		return false
	}

	if contentType := getContentType(filePath); contentType != "" {
		c.Writer.Header().Set("Content-Type", contentType)
	}
	setCacheHeaders(c, filePath)

	if statusCode != http.StatusOK {
		c.Status(statusCode)
		c.Writer.WriteHeaderNow()
		if c.Request.Method != http.MethodHead {
			_, _ = io.Copy(c.Writer, file)
		}
		c.Abort()
		return true
	}

	if seeker, ok := file.(io.ReadSeeker); ok {
		http.ServeContent(c.Writer, c.Request, stat.Name(), stat.ModTime(), seeker)
		c.Abort()
		return true
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return false
	}
	http.ServeContent(c.Writer, c.Request, stat.Name(), stat.ModTime(), bytes.NewReader(data))
	c.Abort()
	return true
}

func getContentType(filePath string) string {
	ext := strings.ToLower(path.Ext(filePath))

	switch ext {
	case ".html", ".htm":
		return "text/html; charset=utf-8"
	case ".css":
		return "text/css; charset=utf-8"
	case ".js":
		return "application/javascript; charset=utf-8"
	case ".json":
		return "application/json; charset=utf-8"
	case ".svg":
		return "image/svg+xml"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".ico":
		return "image/x-icon"
	case ".woff":
		return "font/woff"
	case ".woff2":
		return "font/woff2"
	case ".ttf":
		return "font/ttf"
	default:
		contentType := mime.TypeByExtension(ext)
		if contentType != "" {
			return contentType
		}
		return "application/octet-stream"
	}
}

func setCacheHeaders(c *gin.Context, filePath string) {
	ext := strings.ToLower(path.Ext(filePath))

	if ext == ".html" || ext == ".htm" || ext == ".txt" {
		c.Writer.Header().Set("Cache-Control", "no-cache, must-revalidate")
		c.Writer.Header().Set("Pragma", "no-cache")
		return
	}

	maxAge := 30 * 24 * 60 * 60
	c.Writer.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
	c.Writer.Header().Set("Expires", time.Now().Add(time.Duration(maxAge)*time.Second).Format(time.RFC1123))
}
