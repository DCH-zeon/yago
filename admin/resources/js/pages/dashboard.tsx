import { Head } from '@inertiajs/react';
import AppLayout from '@/layouts/app-layout';
import { home } from '@/routes';
import type { BreadcrumbItem } from '@/types';
import { OnlineUsersWidget } from '@/components/dashboard/online-users-widget';
import { RecentActivityWidget } from '@/components/dashboard/recent-activity-widget';

const breadcrumbs: BreadcrumbItem[] = [
    {
        title: 'Dashboard',
        href: home(),
    },
];

export default function Dashboard() {
    return (
        <AppLayout breadcrumbs={breadcrumbs}>
            <Head title="Dashboard" />
            <div className="flex flex-1 flex-col gap-4 p-4">
                <div className="grid gap-4 md:grid-cols-3">
                    <OnlineUsersWidget />
                    <RecentActivityWidget />
                </div>
            </div>
        </AppLayout>
    );
}
