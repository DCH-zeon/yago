import { Head, useForm, Link } from '@inertiajs/react';
import AppLayout from '@/layouts/app-layout';
import Heading from '@/components/heading';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Label } from '@/components/ui/label';
import type { BreadcrumbItem } from '@/types';
import React from "react";
import InputError from "@/components/input-error";

const breadcrumbs: BreadcrumbItem[] = [
    { title: 'Dashboard', href: '/' },
    { title: 'Доступи', href: '/permissions' },
    { title: 'Створення', href: '/permissions/create' },
];

export default function PermissionCreate() {
    const { data, setData, post, processing, errors } = useForm('post', '/permissions', {
        name: '',
        route: '',
        description: '',
    });

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        post('/permissions', {});
    };

    return (
        <AppLayout breadcrumbs={breadcrumbs}>
            <Head title="Створення доступу" />

            <div className="space-y-4 p-4">
                <Heading
                    title="Створення нового доступу"
                    description="Додайте новий доступ до системи"
                />

                <form onSubmit={handleSubmit} className="max-w-2xl space-y-6">
                    <div className="space-y-2">
                        <div className="flex items-center justify-between mb-1">
                            <Label htmlFor="name">Назва Дії <span className="text-red-500">*</span></Label>
                            <div className="flex items-center gap-2">
                                <span className="text-sm text-muted-foreground hidden sm:block">Застосувати з набору:</span>
                                <Button className="px-2" type="button" variant="outline" size="xs" onClick={() => (setData as any)('name', 'view')}>view</Button>
                                <Button className="px-2" type="button" variant="outline" size="xs" onClick={() => (setData as any)('name', 'create')}>create</Button>
                                <Button className="px-2" type="button" variant="outline" size="xs" onClick={() => (setData as any)('name', 'edit')}>edit</Button>
                                <Button className="px-2" type="button" variant="outline" size="xs" onClick={() => (setData as any)('name', 'delete')}>delete</Button>
                            </div>
                        </div>
                        <Input
                            id="name"
                            value={data.name}
                            onChange={(e) => (setData as any)('name', e.target.value)}
                            placeholder="приклад: view | create | edit | delete"
                            className={errors.name ? 'border-red-500' : ''}
                        />
                        <InputError message={errors.name} />
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="name">Маршрут <span className="text-red-500">*</span></Label>
                        <Input
                            id="route"
                            value={data.route}
                            onChange={(e) => (setData as any)('route', e.target.value)}
                            placeholder="приклад: permissions"
                            className={errors.route ? 'border-red-500' : ''}
                        />
                        <InputError message={errors.route} />
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="description">Опис (опціонально)</Label>
                        <Textarea
                            id="description"
                            value={data.description}
                            onChange={(e) => (setData as any)('description', e.target.value)}
                            placeholder="Описання цього доступу"
                            className={errors.description ? 'border-red-500' : ''}
                        />
                        <InputError message={errors.description} />
                    </div>

                    <div className="flex gap-2">
                        <Button type="submit" disabled={processing}>
                            Створити доступ
                        </Button>
                        <Link href="/permissions">
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
