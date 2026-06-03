import React from "react";

var isInitialized = false;

export function AppWrapper({ children, initialPage }: { children: React.ReactNode, initialPage: any }) {
    const user = initialPage.props.auth.user;

    if (user && !isInitialized && user.settings['is_remember_settings']) {
        const settings = user.settings || {};
        for (let key in settings) {
            if (settings.hasOwnProperty(key)) {
                localStorage.setItem(key, JSON.stringify(settings[key]))
            }
        }
        isInitialized = true;
    }

    return <>{children}</>
}
