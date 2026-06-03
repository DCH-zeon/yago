"use client"

import { ColumnDef } from "@tanstack/react-table";
import { Badge } from '@/components/ui/badge';
import { Link, router } from '@inertiajs/react';
import { Button } from '@/components/ui/button';
import { Pencil, Trash2 } from 'lucide-react';
import {Avatar, AvatarFallback, AvatarImage} from "@/components/ui/avatar";
import React from "react";

type Role = {
    id: number
    name: string
    description: string
}
export type User = {
    id: number;
    name: string;
    email: string;
    created_at: string;
    roles: Role[];
}

const handleDelete = (id: number) => {
    if (confirm('Ви впевнені?')) {
        router.delete(`/users/${id}`);
    }
};

export const getColumns = (isMobile: boolean, permissions: string[]): ColumnDef<User>[] => {

    const canDelete = permissions.includes('delete');
    const canEdit = permissions.includes('edit');

    return [
        {
            id: "Аватар",
            accessorKey: "avatar",
            header: "Аватар",
            cell: ({ row }) => (
                <Avatar className="size-8">
                    {(row.original as any).avatar && (
                        <AvatarImage src={(row.original as any).avatar} alt={row.original.name} />
                    )}
                    <AvatarFallback className="bg-indigo-500 text-white font-semibold">
                        {row.original.name.split(' ').map(n => n[0]).join('').toUpperCase()}
                    </AvatarFallback>
                </Avatar>
            ),
            enableHiding: false,
        },
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
                id: "Пошта",
                accessorKey: "email",
                header: "Пошта",
                cell: ({ row }) => (
                    <div className='text-sm'>
                        {row.original.email}
                    </div>
                ),
            },
            {
                id: "Ролі",
                accessorKey: "roles",
                header: () => <div className="hidden sm:table-cell">Ролі</div>,
                cell: ({ row }) => (
                    <div className="hidden sm:flex flex-wrap gap-1">
                        {
                            row.original.roles.length < 3 ? (
                                row.original.roles.map((role) => (
                                    <Badge key={role.id} variant="secondary" className="px-1.5 text-white text-xs font-normal bg-purple-500 dark:bg-purple-800">
                                        {role.name}
                                    </Badge>))
                            ) : (
                                <Badge variant="secondary" className="px-1.5 text-white text-xs font-normal bg-purple-500 dark:bg-purple-800">
                                    {`${row.original.roles.length} ролі`}
                                </Badge>
                            )
                        }
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
                    const user = row.original
                    return (
                        <div className="flex gap-2 justify-end">
                            {
                                canEdit &&
                                <Link href={`/users/${user.id}/edit`}>
                                    <Button size="sm" variant="outline">
                                        <Pencil className="h-4 w-4" />
                                    </Button>
                                </Link>
                            }
                            {
                                canDelete &&
                                <Button
                                    size="sm"
                                    variant="outline"
                                    onClick={() => handleDelete(user.id)}
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

