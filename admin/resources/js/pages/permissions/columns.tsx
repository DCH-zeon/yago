"use client"

import { ColumnDef } from "@tanstack/react-table";
import { Badge } from '@/components/ui/badge';
import { Link, router } from '@inertiajs/react';
import { Button } from '@/components/ui/button';
import { Pencil, Trash2, SquarePlus, Eye, Route} from 'lucide-react';

export type Permission = {
    id: number
    name: string
    route: string
    guard_name: string
    description: string
    created_at: string
}

const handleDelete = (id: number) => {
    if (confirm('Ви впевнені?')) {
        router.delete(`/permissions/${id}`);
    }
};

export const getColumns = (isMobile: boolean, permissions: string[]): ColumnDef<Permission>[] => {
    const canDelete = permissions.includes('delete');
    const canEdit = permissions.includes('edit');

    return [
        {
            id: "Дозвіл",
            accessorKey: "name",
            header: "Дозвіл",
            cell: ({ row }) => {
                let color;
                let item;
                switch (row.original.name) {
                    case "view":
                        color = "bg-emerald-100 dark:bg-emerald-950";
                        item = <Eye />;
                        break;
                    case "edit":
                        color = "bg-indigo-100 dark:bg-indigo-950";
                        item = <Pencil />;
                        break;
                    case "create":
                        color = "bg-amber-100 dark:bg-yellow-950";
                        item = <SquarePlus />;
                        break;
                    case "delete":
                        color = "bg-red-100 dark:bg-red-950";
                        item = <Trash2 />;
                        break;
                    default:
                        color = "bg-mist-100 dark:bg-mist-800";
                        item = <Route />;
                        break;
                }
                return (
                    <div className="font-medium">
                        <Badge variant="secondary" className={`px-3 ${color}`}>
                            {item}
                            {row.original.name}
                        </Badge>
                    </div>
                )
            },
            enableHiding: false,
        },
        {
            id: "Маршрут",
            accessorKey: "route",
            header: "Маршрут",
            cell: ({ row }) => (
                <div className="text-sm">
                    {row.original.route}
                </div>
            ),
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
                    const permission = row.original
                    return (
                        <div className="flex gap-2 justify-end">
                            {
                                canEdit &&
                                <Link href={`/permissions/${permission.id}/edit`}>
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
                                    onClick={() => handleDelete(permission.id)}
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

export const getFilters = () => [
    {
        id: "Дозвіл",
        accessorKey: "name",
        header: "Дозвіл",
        item: ( value ) => {
            let item;
            switch (value) {
                case "view":
                    item = <Eye />;
                    break;
                case "edit":
                    item = <Pencil />;
                    break;
                case "create":
                    item = <SquarePlus />;
                    break;
                case "delete":
                    item = <Trash2 />;
                    break;
                default:
                    item = <Route />;
                    break;
            }
            return (
                <div className="text-sm flex items-center gap-1">
                        {item} {value}
                </div>
            )
        },
    },
    {
        id: "Маршрут",
        accessorKey: "route",
        header: "Маршрут",
        item: ( value ) => (
            <div className="text-sm">
                {value}
            </div>
        ),
    },
];
