import { createInertiaApp } from '@inertiajs/react';
import { resolvePageComponent } from 'laravel-vite-plugin/inertia-helpers';
import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { TooltipProvider } from '@/components/ui/tooltip';
import '../css/app.css';
import { initializeTheme } from '@/hooks/use-appearance';
import { AppWrapper } from "@/hooks/use-authenticated-init";
import { CentrifugoProvider } from '@/contexts/centrifugo-context';

import { router } from '@inertiajs/react'
const appName = import.meta.env.VITE_APP_NAME || 'Laravel';

createInertiaApp({
    title: (title) => (title ? `${title} - ${appName}` : appName),
    resolve: (name) =>
        resolvePageComponent(
            `./pages/${name}.tsx`,
            import.meta.glob('./pages/**/*.tsx'),
        ),
    setup({el, App, props}) {
        const root = createRoot(el);

        root.render(
            <StrictMode>
                <AppWrapper initialPage={props.initialPage}>
                    <CentrifugoProvider>
                        <TooltipProvider delayDuration={0}>
                            <App {...props} />
                        </TooltipProvider>
                    </CentrifugoProvider>
                </AppWrapper>
            </StrictMode>,
        );
    },

    progress: {
        color: '#4B5563',
    },
})

// This will set light / dark mode on load...
initializeTheme();
