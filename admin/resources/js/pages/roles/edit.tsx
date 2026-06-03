import { Head, useForm, Link } from '@inertiajs/react';
import AppLayout from '@/layouts/app-layout';
import Heading from '@/components/heading';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Checkbox } from '@/components/ui/checkbox';
import { Label } from '@/components/ui/label';
import type { BreadcrumbItem } from '@/types';
import {Textarea} from "@/components/ui/textarea";
import InputError from "@/components/input-error";

interface Permission {
    id: number;
    name: string;
    route: string;
    description?: string;
}

interface Role {
    id: number;
    name: string;
    description?: string;
}

interface RoleFormData {
    name: string;
    description?: string;
    permissions: number[];
}

interface Props {
    role: Role;
    permissions: Permission[];
    rolePermissions: number[];
}

const breadcrumbs: BreadcrumbItem[] = [
    { title: 'Dashboard', href: '/' },
    { title: 'Ролі', href: '/roles' },
    { title: 'Редагування', href: '#' },
];

export default function RoleEdit({ role, permissions, rolePermissions }: Props) {
    const { data, setData, put, processing, errors } = useForm<RoleFormData>('put', `/roles/${role.id}`, {
        name: role.name,
        description: role.description || '',
        permissions: rolePermissions,
    });

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        put(`/roles/${role.id}`,{});
    };

    return (
        <AppLayout breadcrumbs={breadcrumbs}>
            <Head title={`Редагування ролі: ${role.name}`} />

            <div className="space-y-4 p-4">
                <Heading
                    title={`Редагування ролі: ${role.name}`}
                    description="Оновіть назву та доступи ролі"
                />

                <form onSubmit={handleSubmit} className="max-w-2xl space-y-6">
                    <div className="space-y-2">
                        <Label htmlFor="name">Назва ролі  <span className="text-red-500">*</span></Label>
                        <Input
                            id="name"
                            value={data.name}
                            onChange={(e) => (setData as any)('name', e.target.value)}
                            disabled={role.name === 'admin'}
                            className={errors.name ? 'border-red-500' : ''}
                        />
                        <InputError message={errors.name} />
                    </div>

                    <div className="space-y-2">
                        <Label htmlFor="description">Опис (опціонально)</Label>
                        <Textarea
                            id="description"
                            value={data.description}
                            onChange={(e) => (setData as any)('description', e.target.value)}
                            placeholder="Описання ролі"
                            className={errors.description ? 'border-red-500' : ''}
                        />
                        <InputError message={errors.description} />
                    </div>

                    <div className="space-y-3">
                        <Label>Доступи</Label>
                        <div className="space-y-2 max-h-80 overflow-y-auto border rounded p-4">
                            {permissions.map((permission) => (
                                <div key={permission.id} className="flex items-center space-x-2">
                                    <Checkbox
                                        id={`permission-${permission.id}`}
                                        checked={data.permissions.includes(permission.id)}
                                        onCheckedChange={(checked) => {
                                            if (checked) {
                                                (setData as any)('permissions', [
                                                    ...data.permissions,
                                                    permission.id,
                                                ]);
                                            } else {
                                                (setData as any)(
                                                    'permissions',
                                                    data.permissions.filter((id) => id !== permission.id)
                                                );
                                            }
                                        }}
                                    />
                                    <Label htmlFor={`permission-${permission.id}`} className="cursor-pointer">
                                        {permission.description ? permission.description : (permission.name + ' ' + permission.route)}
                                    </Label>
                                </div>
                            ))}
                        </div>
                    </div>

                    <div className="flex gap-2">
                        <Button type="submit" disabled={processing}>
                            Зберегти зміни
                        </Button>
                        <Link href="/roles">
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
