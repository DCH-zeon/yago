import React, { useEffect, useState } from 'react';
import { useCentrifugo } from '@/contexts/centrifugo-context';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {Avatar, AvatarFallback, AvatarImage} from '@/components/ui/avatar';
import { Badge } from '@/components/ui/badge';
import { Users } from 'lucide-react';
import {log} from "node:util";

interface PresenceInfo {
    user: string;
    client: string;
    info?: {
        name?: string;
    };
}

export function OnlineUsersWidget() {
    const { commonSubscription } = useCentrifugo();
    const [users, setUsers] = useState<Record<string, PresenceInfo>>({});

    useEffect(() => {
        // Получаем текущее присутствие при монтировании через API сервера
        fetch('/centrifugo/presence')
            .then(res => res.json())
            .then(data => {
                if (data?.result?.presence) {
                    setUsers(data.result.presence);
                }
            })
            .catch(err => console.error('Failed to fetch presence:', err));

        if (!commonSubscription) return;

        // Слушаем события входа и выхода
        const handleJoin = (ctx: any) => {
            setUsers((prev) => ({
                ...prev,
                [ctx.client]: ctx.info
            }));
        };

        const handleLeave = (ctx: any) => {
            setUsers((prev) => {
                const next = { ...prev };
                delete next[ctx.client];
                return next;
            });
        };

        commonSubscription.on('join', handleJoin);
        commonSubscription.on('leave', handleLeave);

        return () => {
            commonSubscription.removeListener('join', handleJoin);
            commonSubscription.removeListener('leave', handleLeave);
        };
    }, [commonSubscription]);

    const userList = Object.values(users);

    // Удаляем дубликаты пользователей по user ID, чтобы не отображать одного человека несколько раз (разные клиенты/вкладки)
    const uniqueUsers = Object.values(userList.reduce((acc: Record<string, PresenceInfo>, presence) => {
        if (!acc[presence.user] || (presence.info?.name && !acc[presence.user].info?.name)) {
            acc[presence.user] = presence;
        }
        return acc;
    }, {}));
    console.log(uniqueUsers);
    return (
        <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Пользователей онлайн</CardTitle>
                <Users className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
                <div className="text-2xl font-bold mb-4">{uniqueUsers.length}</div>
                <div className="flex flex-wrap gap-2">
                    {uniqueUsers.map((presence) => (
                        <div key={presence.client} className="flex items-center gap-2 bg-secondary p-1 pr-3 rounded-full">
                            <Avatar className="h-6 w-6">
                                <AvatarFallback className="text-[10px]">
                                    {presence.conn_info?.name?.substring(0, 2).toUpperCase() || '??'}
                                </AvatarFallback>
                                <AvatarImage src={presence.conn_info?.avatar} alt={presence.conn_info?.name} className="rounded-full" />
                            </Avatar>
                            <span className="text-xs">{presence.conn_info?.name || 'Unknown'}</span>
                            <Badge variant="outline" className="h-4 px-1 text-[8px] bg-green-500/10 text-green-500 border-green-500/20">
                                online
                            </Badge>
                        </div>
                    ))}
                    {uniqueUsers.length === 0 && (
                        <span className="text-xs text-muted-foreground italic">Никого нет...</span>
                    )}
                </div>
            </CardContent>
        </Card>
    );
}
