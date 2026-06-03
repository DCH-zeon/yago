import { Head, useForm, Link } from '@inertiajs/react';
import AppLayout from '@/layouts/app-layout';
import Heading from '@/components/heading';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Checkbox } from '@/components/ui/checkbox';
import { Label } from '@/components/ui/label';
import type { BreadcrumbItem } from '@/types';
import React, {useState} from "react";
import InputError from "@/components/input-error";
import {Eye, EyeOff} from "lucide-react";

interface Role {
    id: number;
    name: string;
}

interface UserFormData {
    name: string;
    email: string;
    password: string;
    password_confirmation: string;
    roles: number[];
}

interface Props {
    roles: Role[];
}

const breadcrumbs: BreadcrumbItem[] = [
    { title: 'Dashboard', href: '/' },
    { title: 'Користувачі', href: '/users' },
    { title: 'Створення', href: '/users/create' },
];

export default function UserCreate({ roles }: Props) {
    const [showPassword, setShowPassword] = useState(false);
    const { data, setData, post, processing, errors } = useForm<UserFormData>('post', '/users', {
        name: '',
        email: '',
        password: '',
        password_confirmation: '',
        roles: [],
    });

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        post('/users', {});
    };

    return (
        <AppLayout breadcrumbs={breadcrumbs}>
            <Head title="Створення користувача" />

            <div className="space-y-4 p-4">
                <Heading
                    title="Створення нового користувача"
                    description="Додайте нового користувача до системи"
                />

                <form onSubmit={handleSubmit} className="max-w-2xl space-y-6">
                    <div className="space-y-2">
                        <Label htmlFor="name">Ім'я <span className="text-red-500">*</span></Label>
                        <Input
                            id="name"
                            value={data.name}
                            onChange={(e) => (setData as any)('name', e.target.value)}
                            className={errors.name ? 'border-red-500' : ''}
                        />
                        <InputError message={errors.name} />
                    </div>
                    <div className="space-y-2">
                        <Label htmlFor="email">Email <span className="text-red-500">*</span></Label>
                        <Input
                            id="email"
                            type="email"
                            value={data.email}
                            onChange={(e) => (setData as any)('email', e.target.value)}
                            className={errors.email ? 'border-red-500' : ''}
                        />
                        <InputError message={errors.email} />
                    </div>
                    <div className="space-y-2">
                        <Label htmlFor="password">Пароль <span className="text-red-500">*</span></Label>
                        <div className="relative">
                            <Input
                                id="password"
                                type={showPassword ? 'text' : 'password'}
                                value={data.password}
                                onChange={(e) => (setData as any)('password', e.target.value)}
                                placeholder="Новий пароль"
                                className={errors.password ? 'border-red-500' : ''}
                            />
                            <button
                                type="button"
                                onClick={() => setShowPassword(!showPassword)}
                                className="absolute inset-y-0 right-0 pr-3 flex items-center text-sm font-medium text-gray-500 hover:text-gray-700"
                            >
                                {showPassword ? (
                                    <EyeOff className="size-4" />
                                ) : (
                                    <Eye className="size-4" />
                                )}
                            </button>
                        </div>

                        <InputError message={errors.password} />
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="password_confirmation">Підтвердження пароля <span className="text-red-500">*</span></Label>
                        <div className="relative">
                            <Input
                                id="password_confirmation"
                                type={showPassword ? 'text' : 'password'}
                                value={data.password_confirmation}
                                onChange={(e) => (setData as any)('password_confirmation', e.target.value)}
                                className={errors.password_confirmation ? 'border-red-500' : ''}
                            />
                            <button
                                type="button"
                                onClick={() => setShowPassword(!showPassword)}
                                className="absolute inset-y-0 right-0 pr-3 flex items-center text-sm font-medium text-gray-500 hover:text-gray-700"
                            >
                                {showPassword ? (
                                    <EyeOff className="size-4" />
                                ) : (
                                    <Eye className="size-4" />
                                )}
                            </button>
                        </div>

                        <InputError message={errors.password_confirmation} />
                    </div>

                    <div className="space-y-3">
                        <Label>Ролі</Label>
                        <div className="space-y-2 max-h-80 overflow-y-auto border rounded p-4">
                            {roles.map((role) => (
                                <div key={role.id} className="flex items-center space-x-2">
                                    <Checkbox
                                        id={`role-${role.id}`}
                                        checked={data.roles.includes(role.id)}
                                        onCheckedChange={(checked) => {
                                            if (checked) {
                                                (setData as any)('roles', [...data.roles, role.id]);
                                            } else {
                                                (setData as any)(
                                                    'roles',
                                                    data.roles.filter((id) => id !== role.id)
                                                );
                                            }
                                        }}
                                    />
                                    <Label htmlFor={`role-${role.id}`} className="cursor-pointer">
                                        {role.name}
                                    </Label>
                                </div>
                            ))}
                        </div>
                    </div>

                    <div className="flex gap-2">
                        <Button type="submit" disabled={processing}>
                            Створити користувача
                        </Button>
                        <Link href="/users">
                            <Button type="button" variant="outline">
                                Скасувати
                            </Button>
                        </Link>
                    </div>
                </form>
            </div>
        </AppLayout>
    );
}
