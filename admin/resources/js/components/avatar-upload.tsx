import { Upload, X } from 'lucide-react';
import React, { useState } from 'react';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { cn } from '@/lib/utils';

interface AvatarUploadProps {
    currentAvatar?: string;
    userName?: string;
    onFileSelect: (file: File) => void;
    error?: string;
    preview?: string;
    disabled?: boolean;
}

export default function AvatarUpload({
    currentAvatar,
    userName = 'User',
    onFileSelect,
    error,
    preview,
    disabled = false,
}: AvatarUploadProps) {
    const [dragActive, setDragActive] = useState(false);
    const [validationError, setValidationError] = useState<string | null>(null);
    const fileInputRef = React.useRef<HTMLInputElement>(null);

    const handleDrag = (e: React.DragEvent) => {
        e.preventDefault();
        e.stopPropagation();
        if (e.type === 'dragenter' || e.type === 'dragover') {
            setDragActive(true);
        } else if (e.type === 'dragleave') {
            setDragActive(false);
        }
    };

    const handleDrop = (e: React.DragEvent) => {
        e.preventDefault();
        e.stopPropagation();
        setDragActive(false);

        const files = e.dataTransfer.files;
        if (files && files[0]) {
            handleFile(files[0]);
        }
    };

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const files = e.currentTarget.files;
        if (files && files[0]) {
            handleFile(files[0]);
        }
    };

    const handleFile = (file: File) => {
        setValidationError(null);

        // Валідація типу файлу
        const allowedTypes = ['image/jpeg', 'image/png', 'image/jpg', 'image/gif', 'image/webp'];
        if (!allowedTypes.includes(file.type)) {
            setValidationError('Будь ласка, виберіть зображення (JPEG, PNG, GIF, WebP)');
            return;
        }

        // Валідація розміру (макс 5MB)
        if (file.size > 5 * 1024 * 1024) {
            setValidationError('Розмір файлу не повинен перевищувати 5MB');
            return;
        }

        onFileSelect(file);
    };

    const handleButtonClick = () => {
        fileInputRef.current?.click();
    };

    const displayImage = preview || currentAvatar;
    const initials = userName
        .split(' ')
        .map(n => n[0])
        .join('')
        .toUpperCase();

    const displayError = error || validationError;

    return (
        <div className="grid gap-2">
            <label className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                Аватар профілю
            </label>

            <div className="flex flex-col gap-4">
                {/* Превью аватара */}
                <div className="flex items-center gap-4">
                    <Avatar className="size-16">
                        {displayImage && (
                            <AvatarImage src={displayImage} alt={userName} />
                        )}
                        <AvatarFallback>{initials}</AvatarFallback>
                    </Avatar>

                    {/* Зона завантаження */}
                    <div
                        className={cn(
                            'relative flex flex-1 cursor-pointer rounded-lg border-2 border-dashed transition-colors',
                            dragActive
                                ? 'border-primary bg-primary/5'
                                : 'border-muted-foreground/25 bg-muted/10',
                            disabled && 'cursor-not-allowed opacity-50'
                        )}
                        onDragEnter={handleDrag}
                        onDragLeave={handleDrag}
                        onDragOver={handleDrag}
                        onDrop={handleDrop}
                        onClick={handleButtonClick}
                    >
                        <div className="flex w-full items-center justify-center px-4 py-6">
                            <div className="text-center">
                                <Upload className="mx-auto size-5 text-muted-foreground" />
                                <p className="mt-2 text-sm text-muted-foreground">
                                    Перетягніть файл сюди або
                                    <span className="ml-1 font-medium text-foreground">
                                        натисніть для вибору
                                    </span>
                                </p>
                                <p className="mt-1 text-xs text-muted-foreground">
                                    PNG, JPG, GIF до 5MB
                                </p>
                            </div>
                        </div>

                        <input
                            ref={fileInputRef}
                            type="file"
                            name="avatar"
                            accept="image/*"
                            onChange={handleChange}
                            disabled={disabled}
                            className="hidden"
                        />
                    </div>
                </div>

                {/* Повідомлення про помилку */}
                {displayError && (
                    <div className="flex items-start gap-2 rounded-md bg-destructive/10 px-3 py-2">
                        <X className="mt-0.5 size-4 text-destructive" />
                        <p className="text-sm text-destructive">{displayError}</p>
                    </div>
                )}
            </div>
        </div>
    );
}
