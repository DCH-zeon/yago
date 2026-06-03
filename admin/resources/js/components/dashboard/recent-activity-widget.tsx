import React, { useEffect, useState } from 'react';
import { useCentrifugo } from '@/contexts/centrifugo-context';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Activity, Clock } from 'lucide-react';

interface ActivityEvent {
    userId: string | number;
    userName: string;
    action: string;
    url: string;
    timestamp: string;
}

export function RecentActivityWidget() {
    const { commonSubscription } = useCentrifugo();
    const [activities, setActivities] = useState<ActivityEvent[]>([]);
    console.log(commonSubscription)
    useEffect(() => {
        if (!commonSubscription) return;

        const handlePublication = (ctx: any) => {
            console.log('Publication received:', ctx.data);
            const newActivity = ctx.data as ActivityEvent;
            setActivities((prev) => [newActivity, ...prev].slice(0, 10));
        };

        commonSubscription.on('publication', handlePublication);

        return () => {
            commonSubscription.removeListener('publication', handlePublication);
        };
    }, [commonSubscription]);

    return (
        <Card className="col-span-1 md:col-span-2">
            <CardHeader>
                <CardTitle className="text-sm font-medium flex items-center gap-2">
                    <Activity className="h-4 w-4" />
                    Последние действия
                </CardTitle>
            </CardHeader>
            <CardContent>
                <div className="space-y-4">
                    {activities.length > 0 ? (
                        activities.map((activity, idx) => (
                            <div key={`${activity.userId}-${activity.timestamp}-${idx}`} className="flex items-start gap-4 border-b border-sidebar-border/50 pb-3 last:border-0 last:pb-0">
                                <div className="mt-0.5">
                                    <Clock className="h-4 w-4 text-muted-foreground" />
                                </div>
                                <div className="flex-1 space-y-1">
                                    <p className="text-sm font-medium leading-none">
                                        {activity.userName}
                                    </p>
                                    <p className="text-sm text-muted-foreground">
                                        {activity.action}
                                    </p>
                                    <p className="text-[10px] text-muted-foreground opacity-70">
                                        {new Date(activity.timestamp).toLocaleTimeString()} • {new URL(activity.url).pathname}
                                    </p>
                                </div>
                            </div>
                        ))
                    ) : (
                        <div className="text-center py-8 text-muted-foreground text-sm italic">
                            Действий пока не зафиксировано...
                        </div>
                    )}
                </div>
            </CardContent>
        </Card>
    );
}
