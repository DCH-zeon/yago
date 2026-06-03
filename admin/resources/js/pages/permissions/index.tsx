import { Head, Link } from '@inertiajs/react';
import { Plus } from 'lucide-react';
import AppLayout from '@/layouts/app-layout';
import Heading from '@/components/heading';
import { Button } from '@/components/ui/button';
import {getColumns, getFilters} from "./columns"
import { DataTable } from '@/components/data-table'
import type { BreadcrumbItem } from '@/types';
import { useIsMobile } from "@/hooks/use-mobile";
import { useMemo } from "react";

const breadcrumbs: BreadcrumbItem[] = [
    { title: 'Dashboard', href: '/' },
    { title: 'Доступи', href: '/permissions' },
];

export default function PermissionIndex({ auth }) {
    const permissions = auth.permissions['permissions'] ?? [];
    const isMobile = useIsMobile();
    const columns = useMemo(() => getColumns(isMobile, permissions), [isMobile, permissions]);
    const filters = useMemo(() => getFilters(), []);
    const canCreate = permissions && permissions.includes('create');

    return (
        <AppLayout breadcrumbs={breadcrumbs}>
            <Head title="Управління доступами" />
            <div className="space-y-4 p-4">
                <div className="flex items-center justify-between">
                    <Heading
                        title="Доступи"
                        description="Керування доступами системи"
                    />
                    {canCreate && (
                        <Link href="/permissions/create">
                            <Button>
                                <Plus className="h-4 w-4" />
                                Новий доступ
                            </Button>
                        </Link>
                    )}
                </div>
                <DataTable columns={columns} filters={filters} typeKey="permissions"/>
            </div>
        </AppLayout>
    );
}
