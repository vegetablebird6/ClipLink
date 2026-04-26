'use client';

import React, { useState } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faMagnifyingGlass, faSpinner, faXmark } from '@fortawesome/free-solid-svg-icons';

interface SearchBarProps {
  onSearch: (keyword: string) => void | Promise<void>;
  onClear: () => void;
  isSearching?: boolean;
  placeholder?: string;
}

export default function SearchBar({
  onSearch,
  onClear,
  isSearching = false,
  placeholder = '搜索剪贴板内容...'
}: SearchBarProps) {
  const [keyword, setKeyword] = useState('');

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const trimmed = keyword.trim();
    if (trimmed) {
      onSearch(trimmed);
    }
  };

  const handleClear = () => {
    setKeyword('');
    onClear();
  };

  return (
    <form onSubmit={handleSubmit} className="relative w-full">
      <div className="relative">
        <FontAwesomeIcon
          icon={isSearching ? faSpinner : faMagnifyingGlass}
          className={`absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 ${isSearching ? 'animate-spin' : ''}`}
        />
        <input
          value={keyword}
          onChange={(event) => setKeyword(event.target.value)}
          placeholder={placeholder}
          className="w-full rounded-lg border border-gray-200 bg-white py-2 pl-10 pr-10 text-sm text-gray-900 outline-none transition-colors focus:border-blue-500 focus:ring-2 focus:ring-blue-100 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-100 dark:focus:border-blue-400 dark:focus:ring-blue-900/40"
        />
        {keyword && (
          <button
            type="button"
            onClick={handleClear}
            className="absolute right-2 top-1/2 flex h-7 w-7 -translate-y-1/2 items-center justify-center rounded-full text-gray-400 transition-colors hover:bg-gray-100 hover:text-gray-600 dark:hover:bg-gray-700 dark:hover:text-gray-200"
            aria-label="清除搜索"
          >
            <FontAwesomeIcon icon={faXmark} />
          </button>
        )}
      </div>
    </form>
  );
}
