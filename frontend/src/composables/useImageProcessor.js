export const useImageProcessor = () => {
    /**
     * Process an image file to resize it and convert to grayscale
     * @param {File} file - The image file to process
     * @param {Object} options - Processing options
     * @param {number} options.maxSize - Maximum dimension in pixels (default: 300)
     * @param {boolean} options.grayscale - Whether to convert to grayscale (default: true)
     * @param {string} options.format - Output format (default: 'image/jpeg')
     * @param {number} options.quality - Output quality 0-1 (default: 0.85)
     * @returns {Promise<File>} - Processed image file
     */
    const processAvatarImage = (file, options = {}) => {
        const {
            maxSize = 300,
            grayscale = true,
            format = 'image/jpeg',
            quality = 0.85
        } = options;

        return new Promise((resolve, reject) => {
            // Create URL from file
            const url = URL.createObjectURL(file);
            const img = new Image();

            img.onload = () => {
                // Create canvas
                const canvas = document.createElement('canvas');

                // Calculate new dimensions while maintaining aspect ratio
                let width = img.width;
                let height = img.height;

                if (width > height) {
                    if (width > maxSize) {
                        height = Math.round(height * (maxSize / width));
                        width = maxSize;
                    }
                } else {
                    if (height > maxSize) {
                        width = Math.round(width * (maxSize / height));
                        height = maxSize;
                    }
                }

                // Set canvas dimensions
                canvas.width = width;
                canvas.height = height;

                // Draw image to canvas
                const ctx = canvas.getContext('2d');
                ctx.drawImage(img, 0, 0, width, height);

                // Convert to grayscale if requested
                if (grayscale) {
                    const imageData = ctx.getImageData(0, 0, width, height);
                    const data = imageData.data;

                    for (let i = 0; i < data.length; i += 4) {
                        const avg = (data[i] + data[i + 1] + data[i + 2]) / 3;
                        data[i] = avg;     // red
                        data[i + 1] = avg; // green
                        data[i + 2] = avg; // blue
                    }

                    ctx.putImageData(imageData, 0, 0);
                }

                // Convert canvas to blob
                canvas.toBlob((blob) => {
                    // Clean up
                    URL.revokeObjectURL(url);

                    // Create new file from blob
                    const newFile = new File([blob], file.name, {
                        type: format,
                        lastModified: new Date().getTime()
                    });

                    resolve(newFile);
                }, format, quality);
            };

            img.onerror = () => {
                URL.revokeObjectURL(url);
                reject(new Error('Failed to load image'));
            };

            img.src = url;
        });
    };

    /**
     * Get dimensions of an image file
     * @param {File} file - The image file
     * @returns {Promise<{width: number, height: number}>} - Image dimensions
     */
    const getImageDimensions = (file) => {
        return new Promise((resolve, reject) => {
            const url = URL.createObjectURL(file);
            const img = new Image();

            img.onload = () => {
                URL.revokeObjectURL(url);
                resolve({
                    width: img.width,
                    height: img.height
                });
            };

            img.onerror = () => {
                URL.revokeObjectURL(url);
                reject(new Error('Failed to load image'));
            };

            img.src = url;
        });
    };

    return {
        processAvatarImage,
        getImageDimensions
    };
};