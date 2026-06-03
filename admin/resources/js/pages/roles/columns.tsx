"use client"

import { ColumnDef } from "@tanstack/react-table";
import { Badge } from '@/components/ui/badge';
import { Link, router } from '@inertiajs/react';
import { Button } from '@/components/ui/button';
import {Eye, Pencil, Route, SquarePlus, Trash2} from 'lucide-react';

export type Role = {
    id: number
    name: string
    guard_name: string
    description: string
    created_at: string
    permissions: Array<{ id: number; name: string }>;
}

const handleDelete = (id: number, name: string) => {
    if (name === 'admin') {
        alert('Неможливо видалити роль адміністратора');
        return;
    }
    if (confirm('Ви впевнені?')) {
        router.delete(`/roles/${id}`);
    }
};

export const getColumns = (isMobile: boolean, permissions: string[]): ColumnDef<Role>[] => {
    const canDelete = permissions.includes('delete');
    const canEdit = permissions.includes('edit');

    return [
        {
            id: "Назва",
            accessorKey: "name",
            header: "Назва",
            cell: ({ row }) => (
                <div className="font-medium">
                    {row.original.name}
                </div>
            ),
            enableHiding: false,
        },
        ...(isMobile ? [] : [
            {
                id: "Опис",
                accessorKey: "description",
                header: "Опис",
                cell: ({ row }) => (
                    <div className={row.original.description ? 'text-sm' : 'text-sm text-muted-foreground italic'}>
                        {row.original.description ? (row.original.description.length > 50 ? row.original.description.slice(0, 50) + '...' : row.original.description) : 'Немає опису'}
                    </div>
                ),
            },
            {
                id: "Тип захисту",
                accessorKey: "guard_name",
                header: () => <div className="w-32 hidden sm:table-cell">Тип захисту</div>,
                cell: ({ row }) => (
                    <div className="w-32 hidden sm:table-cell">
                        <Badge variant="outline" className="px-1.5 text-muted-foreground">
                            {row.original.guard_name}
                        </Badge>
                    </div>
                ),
            },
            {
                id: "Доступи",
                accessorKey: "permissions",
                header: () => <div className="w-32 hidden sm:table-cell">Доступи</div>,
                cell: ({ row }) => (
                    <div className="w-32 hidden sm:table-cell">
                        {row.original.permissions.length} доступів
                    </div>
                ),
            },
            {
                id: "Дата створення",
                accessorKey: "created_at",
                header: () => <div className="w-32 hidden sm:table-cell">Дата створення</div>,
                cell: ({ row }) => {
                    const date = new Date(row.original.created_at);
                    const formatted = date.toLocaleDateString('uk-UA', {
                        year: 'numeric',
                        month: '2-digit',
                        day: '2-digit'
                    });
                    return (
                        <div className="w-32 hidden sm:table-cell">
                            {formatted}
                        </div>
                    );
                },
            },
        ]),
        ...(canEdit || canDelete ? [
            {
                id: "actions",
                header: () => <div className="flex justify-end mr-10">Дії</div>,
                cell: ({ row }) => {
                    const role = row.original
                    return (
                        <div className="flex gap-2 justify-end">
                            {
                                canEdit &&
                                <Link href={`/roles/${role.id}/edit`}>
                                    <Button size="sm" variant="outline">
                                        <Pencil className="h-4 w-4" />
                                    </Button>
                                </Link>
                            }
                            {
                                role.name !== 'admin' && canDelete &&
                                <Button
                                    size="sm"
                                    variant="outline"
                                    onClick={() => handleDelete(role.id, role.name)}
                                >
                                    <Trash2 className="h-4 w-4 text-red-500" />
                                </Button>
                            }
                        </div>
                    )
                },
                enableSorting: false,
                enableHiding: false,
            },
        ] : [])
    ]
}

export const getFilters = () => [];

