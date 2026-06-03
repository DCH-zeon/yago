import { Head, Link, router } from '@inertiajs/react';
import { Plus, Pencil, Trash2, Shield } from 'lucide-react';
import AppLayout from '@/layouts/app-layout';
import Heading from '@/components/heading';
import { Button } from '@/components/ui/button';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import type { BreadcrumbItem } from '@/types';
import { getColumns, User } from "./columns";
import {useIsMobile} from "@/hooks/use-mobile";
import {useMemo} from "react";
import {DataTable} from "@/components/data-table";

const breadcrumbs: BreadcrumbItem[] = [
    { title: 'Dashboard', href: '/' },
    { title: 'Користувачі', href: '/users' },
];

export default function UserIndex({ auth }) {
    const permissions = auth.permissions['permissions'] ?? [];
    const isMobile = useIsMobile();
    const columns = useMemo(() => getColumns(isMobile, permissions), [isMobile, permissions]);
    const filters = [];
    const canCreate = permissions && permissions.includes('create');



    const typeKey = 'users';

    return (
        <AppLayout breadcrumbs={breadcrumbs}>
            <Head title="Управління користувачами" />

            <div className="space-y-4 p-4">
                <div className="flex items-center justify-between">
                    <Heading
                        title="Користувачі"
                        description="Керування користувачами системи"
                    />
                    {canCreate && (
                        <Link href="/users/create">
                            <Button>
                                <Plus className="h-4 w-4" />
                                Новий користувач
                            </Button>
                        </Link>
                    )}
                </div>
                <DataTable columns={columns} filters={filters} typeKey="users"/>
                {/*<div className="border rounded-lg overflow-hidden">*/}
                {/*    /!*<Table>*!/*/}
                {/*    /!*    <TableHeader>*!/*/}
                {/*    /!*        <TableRow className="bg-muted/50">*!/*/}
                {/*    /!*            <TableHead>Ім'я</TableHead>*!/*/}
                {/*    /!*            <TableHead>Email</TableHead>*!/*/}
                {/*    /!*            <TableHead>Ролі</TableHead>*!/*/}
                {/*    /!*            <TableHead className="w-40">Дії</TableHead>*!/*/}
                {/*    /!*        </TableRow>*!/*/}
                {/*    /!*    </TableHeader>*!/*/}
                {/*    /!*    <TableBody>*!/*/}
                {/*    /!*        {users.data.map((user) => (*!/*/}
                {/*    /!*            <TableRow key={user.id}>*!/*/}
                {/*    /!*                <TableCell className="font-medium">{user.name}</TableCell>*!/*/}
                {/*    /!*                <TableCell className="text-sm">{user.email}</TableCell>*!/*/}
                {/*    /!*                <TableCell>*!/*/}
                {/*    /!*                    <div className="flex gap-1 flex-wrap">*!/*/}
                {/*    /!*                        {user.roles.length > 0 ? (*!/*/}
                {/*    /!*                            user.roles.map((role) => (*!/*/}
                {/*    /!*                                <Badge key={role.id} variant="secondary">*!/*/}
                {/*    /!*                                    {role.name}*!/*/}
                {/*    /!*                                </Badge>*!/*/}
                {/*    /!*                            ))*!/*/}
                {/*    /!*                        ) : (*!/*/}
                {/*    /!*                            <span className="text-gray-500 text-sm">Немає ролей</span>*!/*/}
                {/*    /!*                        )}*!/*/}
                {/*    /!*                    </div>*!/*/}
                {/*    /!*                </TableCell>*!/*/}
                {/*    /!*                <TableCell className="flex gap-2">*!/*/}
                {/*    /!*                    <Link href={`/users/${user.id}/edit`}>*!/*/}
                {/*    /!*                        <Button size="sm" variant="outline">*!/*/}
                {/*    /!*                            <Pencil className="h-4 w-4" />*!/*/}
                {/*    /!*                        </Button>*!/*/}
                {/*    /!*                    </Link>*!/*/}
                {/*    /!*                    <Link href={`/users/${user.id}/edit`}>*!/*/}
                {/*    /!*                        <Button size="sm" variant="outline" title="Налаштування ролей">*!/*/}
                {/*    /!*                            <Shield className="h-4 w-4" />*!/*/}
                {/*    /!*                        </Button>*!/*/}
                {/*    /!*                    </Link>*!/*/}
                {/*    /!*                    <Button*!/*/}
                {/*    /!*                        size="sm"*!/*/}
                {/*    /!*                        variant="outline"*!/*/}
                {/*    /!*                        onClick={() => handleDelete(user.id)}*!/*/}
                {/*    /!*                    >*!/*/}
                {/*    /!*                        <Trash2 className="h-4 w-4 text-red-500" />*!/*/}
                {/*    /!*                    </Button>*!/*/}
                {/*    /!*                </TableCell>*!/*/}
                {/*    /!*            </TableRow>*!/*/}
                {/*    /!*        ))}*!/*/}
                {/*    /!*    </TableBody>*!/*/}
                {/*    /!*</Table>*!/*/}
                {/*</div>*/}
            </div>
        </AppLayout>
    );
}
