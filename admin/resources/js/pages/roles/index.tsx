import { Head, Link } from '@inertiajs/react';
import { Plus } from 'lucide-react';
import AppLayout from '@/layouts/app-layout';
import Heading from '@/components/heading';
import { Button } from '@/components/ui/button';
import type { BreadcrumbItem } from '@/types';
import { getColumns, getFilters } from "./columns"
import {useIsMobile} from "@/hooks/use-mobile";
import {useMemo} from "react";
import {DataTable} from "@/components/data-table";

const breadcrumbs: BreadcrumbItem[] = [
    { title: 'Dashboard', href: '/' },
    { title: 'Ролі', href: '/roles' },
];

export default function RoleIndex({ auth }) {
    const permissions = auth.permissions['permissions'] ?? [];
    const isMobile = useIsMobile();
    const columns = useMemo(() => getColumns(isMobile, permissions), [isMobile, permissions]);
    const filters = [];
    const canCreate = permissions && permissions.includes('create');

    return (
        <AppLayout breadcrumbs={breadcrumbs}>
            <Head title="Управління ролями" />

            <div className="space-y-4 p-4">
                <div className="flex items-center justify-between">
                    <Heading
                        title="Ролі"
                        description="Керування ролями користувачів системи"
                    />
                    {canCreate && (
                        <Link href="/roles/create">
                            <Button>
                                <Plus className="h-4 w-4" />
                                Нова роль
                            </Button>
                        </Link>
                    )}
                </div>
                <DataTable columns={columns} filters={filters} typeKey="roles"/>
            </div>
        </AppLayout>
    );
}
