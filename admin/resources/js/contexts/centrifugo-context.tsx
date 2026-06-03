import React, { createContext, useContext, useEffect, useState, useRef } from 'react';
import { Centrifuge, Subscription } from 'centrifuge';

interface CentrifugoContextType {
    centrifuge: Centrifuge | null;
    commonSubscription: Subscription | null;
}

const CentrifugoContext = createContext<CentrifugoContextType>({
    centrifuge: null,
    commonSubscription: null,
});

export const useCentrifugo = () => useContext(CentrifugoContext);

export const CentrifugoProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [centrifuge, setCentrifuge] = useState<Centrifuge | null>(null);
    const [commonSubscription, setCommonSubscription] = useState<Subscription | null>(null);
    const centrifugeRef = useRef<Centrifuge | null>(null);

    useEffect(() => {
        const initCentrifuge = async () => {
            try {
                const response = await fetch('/centrifugo/token');
                const { token } = await response.json();

                // Використовуємо URL із налаштувань Laravel (через shared data) або визначаємо автоматично
                let centrifugeUrl = (window as any).laravel?.centrifugo_url;

                if (!centrifugeUrl) {
                    // Намагаємося взяти з мета-тегу (якщо передаємо через нього) або визначаємо по поточному хосту
                    centrifugeUrl = document.querySelector('meta[name="centrifugo-url"]')?.getAttribute('content');
                }

                if (!centrifugeUrl) {
                    const host = window.location.hostname;
                    const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
                    centrifugeUrl = `${protocol}://${host}/connection/websocket`;
                }

                const cf = new Centrifuge(centrifugeUrl, {
                    token: token,
                });

                cf.on('connected', (ctx) => {
                    console.log('Connected to Centrifugo', ctx);
                });

                cf.on('error', (ctx) => {
                    console.error('Centrifugo error', ctx);
                });

                cf.on('disconnected', (ctx) => {
                    console.log('Disconnected from Centrifugo', ctx);
                });

                cf.connect();
                centrifugeRef.current = cf;
                setCentrifuge(cf);

                const sub = cf.newSubscription('admin');

                // Уключаємо отримання даних про присутність (presence) та події join/leave
                // У Centrifugo це також має бути дозволено в налаштуваннях каналу
                sub.on('publication', (ctx) => {
                    console.log('Received publication', ctx);
                });

                sub.on('join', (ctx) => {
                    console.log('User joined', ctx);
                });

                sub.on('leave', (ctx) => {
                    console.log('User left', ctx);
                });

                sub.subscribe();
                setCommonSubscription(sub);
            } catch (error) {
                console.error('Centrifugo connection error:', error);
            }
        };

        initCentrifuge();

        return () => {
            if (centrifugeRef.current) {
                centrifugeRef.current.disconnect();
            }
        };
    }, []);

    return (
        <CentrifugoContext.Provider value={{ centrifuge, commonSubscription }}>
            {children}
        </CentrifugoContext.Provider>
    );
};
