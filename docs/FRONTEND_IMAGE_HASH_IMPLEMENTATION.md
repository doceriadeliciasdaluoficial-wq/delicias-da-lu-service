# Frontend Image Implementation Guide

## Overview

This guide explains how to implement image handling in the Delícias da Lú frontend. Images are stored as base64-encoded strings directly in Firestore object fields:

1. **Upload**: Frontend loads image bytes → converts to base64 string → sends to backend
2. **Storage**: Backend stores base64 string directly in the `Image` field of Firestore documents
3. **Retrieve**: Backend returns complete object including the base64 `Image` field
4. **Display**: Frontend converts base64 string back to blob/URL for display
5. **Important**: Frontend does NOT store image data locally

## Architecture

```
┌─────────────┐                    ┌──────────────┐
│   Frontend  │                    │   Backend    │
│  (React/TS) │                    │   (Go/Echo)  │
└─────┬───────┘                    └──────┬───────┘
      │                                   │
      ├─► Load image from file/URL        │
      ├─► Read image bytes                │
      ├─► Convert to base64 string        │
      │                                   │
      ├─────── Send: {Image: "base64..."}→│
      │                                   │
      │                          Store in Firestore:
      │                            Image field contains
      │                            base64 string
      │                                   │
      │     ◄───── Returns: object with Image field
      │                                   │
      └─► Convert base64 to blob         │
          Display as <img src="blob:..."/> │
```

## Implementation

### 1. Types and Interfaces

```typescript
// types/image.ts

/**
 * Configuration for image operations
 */
export interface ImageConfig {
  /** Max file size in MB (default: 5) */
  maxSizeMB?: number;
  
  /** Allowed MIME types (default: jpeg, png, webp) */
  allowedMimes?: string[];
}

/**
 * Image upload payload sent to backend
 */
export interface ImageUploadPayload {
  /** Base64-encoded image string */
  image: string;
}

/**
 * MenuItem with image base64
 */
export interface MenuItemWithImage extends MenuItem {
  image?: string; // base64-encoded image
}

/**
 * CakeBuilder component with image base64
 */
export interface CakeBuilderComponentWithImage extends CakeBuilderComponent {
  image?: string; // base64-encoded image
}

/**
 * Home/Config object with image base64
 */
export interface HomeWithImage extends Home {
  image?: string; // base64-encoded image
}
```

### 2. Image Utility Functions

```typescript
// utils/imageUtils.ts

/**
 * Convert Blob/File to base64 string
 */
export async function fileToBase64(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => {
      const result = reader.result as string;
      // Remove 'data:image/...;base64,' prefix
      const base64 = result.split(',')[1];
      resolve(base64);
    };
    reader.onerror = reject;
    reader.readAsDataURL(file);
  });
}

/**
 * Convert base64 string back to Blob
 */
export function base64ToBlob(
  base64: string,
  mimeType: string = 'image/jpeg'
): Blob {
  const binaryString = atob(base64);
  const bytes = new Uint8Array(binaryString.length);
  
  for (let i = 0; i < binaryString.length; i++) {
    bytes[i] = binaryString.charCodeAt(i);
  }
  
  return new Blob([bytes], { type: mimeType });
}

/**
 * Create object URL from base64 string for displaying in <img>
 */
export function base64ToObjectURL(
  base64: string,
  mimeType: string = 'image/jpeg'
): string {
  const blob = base64ToBlob(base64, mimeType);
  return URL.createObjectURL(blob);
}

/**
 * Load image from URL and convert to base64
 */
export async function loadImageFromURL(url: string): Promise<string> {
  const response = await fetch(url);
  
  if (!response.ok) {
    throw new Error(`Failed to load image from URL: ${response.statusText}`);
  }

  const blob = await response.blob();
  return fileToBase64(blob as File);
}

/**
 * Validate image file
 */
export async function validateImageFile(
  file: File,
  maxSizeMB: number = 5,
  allowedMimes: string[] = ['image/jpeg', 'image/png', 'image/webp']
): Promise<{ valid: boolean; error?: string }> {
  // Check file size
  if (file.size > maxSizeMB * 1024 * 1024) {
    return {
      valid: false,
      error: `File size exceeds ${maxSizeMB}MB limit`,
    };
  }

  // Check MIME type
  if (!allowedMimes.includes(file.type)) {
    return {
      valid: false,
      error: `File type ${file.type} not allowed. Allowed: ${allowedMimes.join(', ')}`,
    };
  }

  return { valid: true };
}
```

  return { valid: true };
}

/**
 * Load image from URL
 */
export async function loadImageFromURL(
  url: string
): Promise<{ blob: Blob; mimeType: string }> {
  const response = await fetch(url);
  
  if (!response.ok) {
    throw new Error(`Failed to load image from URL: ${response.statusText}`);
  }

  const blob = await response.blob();
  const mimeType = response.headers.get('content-type') || 'image/jpeg';

  return { blob, mimeType };
}

/**
 * Load image from file input
 */
export async function loadImageFromFileInput(
  file: File
): Promise<{ blob: Blob; mimeType: string; filename: string }> {
  return {
    blob: file,
    mimeType: file.type,
    filename: file.name,
  };
}
```

### 3. Image API Service

```typescript
// services/imageService.ts

import { fileToBase64, validateImageFile } from '@/utils/imageUtils';
import { ImageConfig } from '@/types/image';

/**
 * Service for handling image uploads and API communication
 */
export class ImageService {
  private config: ImageConfig;

  constructor(config: Partial<ImageConfig> = {}) {
    this.config = {
      maxSizeMB: 5,
      allowedMimes: ['image/jpeg', 'image/png', 'image/webp'],
      apiBaseUrl: process.env.REACT_APP_API_URL || 'http://localhost:8080',
      ...config,
    };
  }

  /**
   * Convert file to base64 and validate
   */
  async processImageFile(file: File): Promise<string> {
    // Validate file
    const validation = await validateImageFile(
      file,
      this.config.maxSizeMB,
      this.config.allowedMimes
    );

    if (!validation.valid) {
      throw new Error(validation.error);
    }

    // Convert to base64
    return fileToBase64(file);
  }

  /**
   * Create menu item with image
   */
  async createMenuItem(data: any): Promise<any> {
    const response = await fetch(`${this.config.apiBaseUrl}/api/menu`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      throw new Error(`Failed to create menu item: ${response.statusText}`);
    }

    return response.json();
  }

  /**
   * Update menu item with image
   */
  async updateMenuItem(id: string, data: any): Promise<any> {
    const response = await fetch(`${this.config.apiBaseUrl}/api/menu/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      throw new Error(`Failed to update menu item: ${response.statusText}`);
    }

    return response.json();
  }

  /**
   * Get menu items (includes base64 images in response)
   */
  async getMenuItems(): Promise<any[]> {
    const response = await fetch(`${this.config.apiBaseUrl}/api/menu`);

    if (!response.ok) {
      throw new Error(`Failed to fetch menu items: ${response.statusText}`);
    }

    return response.json();
  }
}

export const imageService = new ImageService();
```

### 4. React Hooks

```typescript
// hooks/useImageUpload.ts

import { useState, useCallback } from 'react';
import React from 'react';
import { imageService } from '@/services/imageService';
import { base64ToObjectURL } from '@/utils/imageUtils';

/**
 * Hook for image upload and conversion
 */
export function useImageUpload() {
  const [base64, setBase64] = useState<string | null>(null);
  const [displayURL, setDisplayURL] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [mimeType, setMimeType] = useState<string>('image/jpeg');

  const handleFileUpload = useCallback(async (file: File) => {
    setLoading(true);
    setError(null);

    try {
      // Process and convert file to base64
      const b64 = await imageService.processImageFile(file);
      setBase64(b64);
      setMimeType(file.type);

      // Create display URL
      const url = base64ToObjectURL(b64, file.type);
      setDisplayURL(url);

      return b64;
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Upload failed';
      setError(message);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  const handleURLUpload = useCallback(async (url: string) => {
    setLoading(true);
    setError(null);

    try {
      const response = await fetch(url);
      const blob = await response.blob();
      const file = new File([blob], 'image', { type: blob.type });
      
      return handleFileUpload(file);
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to load image';
      setError(message);
      throw err;
    } finally {
      setLoading(false);
    }
  }, [handleFileUpload]);

  const reset = useCallback(() => {
    setBase64(null);
    setDisplayURL(null);
    setError(null);
    if (displayURL?.startsWith('blob:')) {
      URL.revokeObjectURL(displayURL);
    }
  }, [displayURL]);

  return {
    base64,
    displayURL,
    mimeType,
    loading,
    error,
    handleFileUpload,
    handleURLUpload,
    reset,
  };
}

/**
 * Hook to display base64 image from server
 */
export function useImageDisplay(imageBase64?: string) {
  const [displayURL, setDisplayURL] = useState<string | null>(null);

  React.useEffect(() => {
    if (!imageBase64) {
      setDisplayURL(null);
      return;
    }

    try {
      // Assume JPEG by default, but could be passed as parameter
      const url = base64ToObjectURL(imageBase64, 'image/jpeg');
      setDisplayURL(url);

      // Cleanup on unmount
      return () => {
        if (url.startsWith('blob:')) {
          URL.revokeObjectURL(url);
        }
      };
    } catch (err) {
      console.error('Failed to display image:', err);
      setDisplayURL(null);
    }
  }, [imageBase64]);

  return displayURL;
}
```

### 5. React Components

```typescript
// components/ImageUploader.tsx

import React, { useRef } from 'react';
import { useImageUpload } from '@/hooks/useImageUpload';

interface ImageUploaderProps {
  onImageBase64Change: (base64: string | null, mimeType: string) => void;
  label?: string;
  initialImage?: string;
}

export const ImageUploader: React.FC<ImageUploaderProps> = ({
  onImageBase64Change,
  label = 'Upload Image',
  initialImage,
}) => {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const {
    base64,
    displayURL,
    mimeType,
    loading,
    error,
    handleFileUpload,
    handleURLUpload,
  } = useImageUpload();

  React.useEffect(() => {
    onImageBase64Change(base64, mimeType);
  }, [base64, mimeType, onImageBase64Change]);

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      handleFileUpload(file).catch((err) => console.error(err));
    }
  };

  const handleURLChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const url = event.target.value.trim();
    if (url) {
      handleURLUpload(url).catch((err) => console.error(err));
    }
  };

  return (
    <div className="image-uploader">
      <label>{label}</label>

      {/* File input */}
      <input
        type="file"
        ref={fileInputRef}
        onChange={handleFileChange}
        accept="image/jpeg,image/png,image/webp"
        disabled={loading}
        hidden
      />
      <button onClick={() => fileInputRef.current?.click()} disabled={loading}>
        {loading ? 'Processing...' : 'Choose File'}
      </button>

      {/* URL input */}
      <input
        type="url"
        placeholder="Or paste image URL"
        onChange={handleURLChange}
        disabled={loading}
      />

      {/* Preview */}
      {displayURL && (
        <div className="image-preview">
          <img src={displayURL} alt="Preview" />
          <p className="image-size">
            {base64 ? `${(base64.length * 0.75 / 1024).toFixed(1)} KB` : ''}
          </p>
        </div>
      )}

      {/* Error message */}
      {error && <div className="error-message">{error}</div>}
    </div>
  );
};

export default ImageUploader;
```

```typescript
// components/ImageDisplay.tsx

import React from 'react';
import { base64ToObjectURL } from '@/utils/imageUtils';
import { useImageDisplay } from '@/hooks/useImageUpload';

interface ImageDisplayProps {
  imageBase64?: string;
  alt?: string;
  className?: string;
  mimeType?: string;
}

export const ImageDisplay: React.FC<ImageDisplayProps> = ({
  imageBase64,
  alt = 'Image',
  className,
  mimeType = 'image/jpeg',
}) => {
  const [displayURL, setDisplayURL] = React.useState<string | null>(null);

  React.useEffect(() => {
    if (!imageBase64) {
      setDisplayURL(null);
      return;
    }

    try {
      const url = base64ToObjectURL(imageBase64, mimeType);
      setDisplayURL(url);

      return () => {
        if (url.startsWith('blob:')) {
          URL.revokeObjectURL(url);
        }
      };
    } catch (err) {
      console.error('Failed to display image:', err);
      setDisplayURL(null);
    }
  }, [imageBase64, mimeType]);

  if (!imageBase64) {
    return <div className={className}>No image</div>;
  }

  if (displayURL) {
    return <img src={displayURL} alt={alt} className={className} />;
  }

  return <div className={className}>Unable to display image</div>;
};

export default ImageDisplay;
```

### 6. Usage in Forms

```typescript
// components/MenuItemForm.tsx

import React, { useState } from 'react';
import ImageUploader from '@/components/ImageUploader';
import { imageService } from '@/services/imageService';

export const MenuItemForm: React.FC = () => {
  const [formData, setFormData] = useState({
    name: '',
    category: '',
    price: 0,
    image: '', // base64 string
    description: '',
    active: true,
    order: 0,
  });

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  const handleImageChange = (base64: string | null, mimeType: string) => {
    setFormData((prev) => ({
      ...prev,
      image: base64 || '',
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!formData.image) {
      setError('Please upload an image');
      return;
    }

    setLoading(true);
    setError(null);
    setSuccess(false);

    try {
      // Send to backend with base64 image
      await imageService.createMenuItem({
        ...formData,
        image: formData.image, // Base64 string is stored directly
      });

      setSuccess(true);
      setFormData({
        name: '',
        category: '',
        price: 0,
        image: '',
        description: '',
        active: true,
        order: 0,
      });
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Failed to create item';
      setError(message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="menu-item-form">
      <input
        type="text"
        placeholder="Item name"
        value={formData.name}
        onChange={(e) =>
          setFormData((prev) => ({ ...prev, name: e.target.value }))
        }
        required
      />

      <input
        type="text"
        placeholder="Category"
        value={formData.category}
        onChange={(e) =>
          setFormData((prev) => ({ ...prev, category: e.target.value }))
        }
        required
      />

      <input
        type="number"
        placeholder="Price"
        value={formData.price}
        onChange={(e) =>
          setFormData((prev) => ({ ...prev, price: parseFloat(e.target.value) }))
        }
        required
      />

      <textarea
        placeholder="Description"
        value={formData.description}
        onChange={(e) =>
          setFormData((prev) => ({ ...prev, description: e.target.value }))
        }
      />

      <ImageUploader
        onImageBase64Change={handleImageChange}
        label="Item Image"
      />

      {error && <div className="error">{error}</div>}
      {success && <div className="success">Created successfully!</div>}

      <button type="submit" disabled={loading}>
        {loading ? 'Creating...' : 'Create'}
      </button>
    </form>
  );
};

export default MenuItemForm;
```

### 7. Retrieving and Displaying Stored Items

```typescript
// components/MenuItemList.tsx

import React, { useEffect, useState } from 'react';
import ImageDisplay from '@/components/ImageDisplay';
import { imageService } from '@/services/imageService';

interface MenuItem {
  id: string;
  name: string;
  category: string;
  price: number;
  image?: string; // base64 string from backend
  description?: string;
  active: boolean;
}

export const MenuItemList: React.FC = () => {
  const [items, setItems] = useState<MenuItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchItems = async () => {
      try {
        // Backend returns items with base64 images in the Image field
        const data = await imageService.getMenuItems();
        setItems(data);
      } catch (err) {
        const message = err instanceof Error ? err.message : 'Failed to fetch items';
        setError(message);
      } finally {
        setLoading(false);
      }
    };

    fetchItems();
  }, []);

  if (loading) return <div className="loading">Loading menu items...</div>;
  if (error) return <div className="error">{error}</div>;

  return (
    <div className="menu-list">
      {items.map((item) => (
        <div key={item.id} className="menu-item">
          <h3>{item.name}</h3>
          
          {/* Display image from base64 */}
          {item.image ? (
            <ImageDisplay
              imageBase64={item.image}
              alt={item.name}
              className="menu-item-image"
              mimeType="image/jpeg"
            />
          ) : (
            <div className="no-image">No image</div>
          )}
          
          <p className="category">{item.category}</p>
          <p className="price">${item.price.toFixed(2)}</p>
          <p className="description">{item.description}</p>
        </div>
      ))}
    </div>
  );
};

export default MenuItemList;
```

## Key Points

### Data Flow

1. **Upload**: User selects file → converted to base64 → sent to backend
2. **Storage**: Backend receives base64 → stores in Firestore Image field
3. **Retrieval**: Backend returns object with Image field (base64)
4. **Display**: Frontend converts base64 → blob URL → displays in &lt;img&gt;

### Memory Management

```typescript
// Always clean up object URLs to prevent memory leaks
if (displayURL && displayURL.startsWith('blob:')) {
  URL.revokeObjectURL(displayURL);
}
```

### Base64 Size Note

- Base64 encoding increases size by ~33%
- A 5MB image becomes ~6.7MB in base64
- Adjust max file size if needed

### Security Considerations

1. **Validate files** on both client and server sides
2. **File size limits** - Enforce maximum file sizes (5MB recommended)
3. **MIME type validation** - Only allow jpeg, png, webp
4. **HTTPS only** - Use HTTPS in production
5. **Input sanitization** - Validate URLs before fetching

### Performance Tips

1. **Lazy load images** - Only convert base64 to blob URL when displaying
2. **Cleanup object URLs** - Revoke blob URLs after use to free memory
3. **Progressive upload** - Show preview before sending to backend
4. **Error handling** - Gracefully handle upload/fetch failures
5. **WebP support** - Use modern formats for smaller file sizes

## Backend Integration

When sending data to the backend:

```typescript
// Create menu item with image hash
const createMenuItemRequest = {
  id: 'item-1',
  name: 'Chocolate Cake',
  category: 'Bolos',
  price: 45.00,
  image: 'a1b2c3d4e5f6...', // SHA-256 hash (64 hex characters)
  description: 'Delicious chocolate cake',
  active: true,
};

// Backend stores the hash string in Firestore:
// {
//   id: 'item-1',
//   name: 'Chocolate Cake',
//   image: 'a1b2c3d4e5f6...',  // Just the hash
//   ...
// }
```

When retrieving:

```typescript
// Backend returns object with hash
const menuItem = await api.get('/menu/items/1');
// Returns: { id: '1', name: 'Chocolate Cake', image: 'a1b2c3d4e5f6...', ... }

// Frontend uses hash to display
<ImageDisplay hash={menuItem.image} alt={menuItem.name} />
```

## Migration Path

If migrating from URL-based images:

1. Create a migration script to hash all existing images
2. Update image fields in Firestore from URLs to hashes
3. Deploy backend changes first
4. Deploy frontend changes
5. Run migration on existing data
6. Remove old URL references

## Dependencies

Install required packages:

```bash
npm install crypto-js
# or use Web Crypto API (no installation needed for modern browsers)
```

For TypeScript types:

```bash
npm install --save-dev @types/crypto-js
```

## Troubleshooting

### Images not displaying
- Check browser console for errors
- Verify hash exists in storage
- Ensure image data wasn't corrupted

### Memory leaks
- Always revoke object URLs when done
- Use cleanup functions in useEffect

### Storage quota exceeded
- Switch to IndexedDB (larger quota)
- Implement image cleanup/deletion
- Show user warning before quota is exceeded

### Hash mismatch
- Ensure same algorithm (SHA-256) used everywhere
- Verify image data isn't corrupted during transfer
- Check for encoding issues (base64)

## References

- [Web Crypto API](https://developer.mozilla.org/en-US/docs/Web/API/Web_Crypto_API)
- [IndexedDB](https://developer.mozilla.org/en-US/docs/Web/API/IndexedDB_API)
- [LocalStorage](https://developer.mozilla.org/en-US/docs/Web/API/Window/localStorage)
- [Blob API](https://developer.mozilla.org/en-US/docs/Web/API/Blob)
