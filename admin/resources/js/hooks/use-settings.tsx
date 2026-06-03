const setCookie = (name: string, value: string, days = 365): void => {
    if (typeof document === 'undefined') {
        return;
    }

    const maxAge = days * 24 * 60 * 60;
    document.cookie = `${name}=${encodeURIComponent(value)};path=/;max-age=${maxAge};SameSite=Lax`;
};

const getCookie = (name: string) => {
    if (typeof document === 'undefined') {
        return;
    }

    const cookies = document.cookie.split("; ");
    for (let c of cookies) {
        const [key, value] = c.split("=");
        if (key === name) {
            return decodeURIComponent(value);
        }
    }
    return null;
}

const setLocalStorage = (name: string, value: string): void => {
    if (typeof window === 'undefined') {
        return;
    }
    try {
        localStorage.setItem(name, value);
    } catch (error) {}
};

const deleteLocalStorage = (name: string): void => {
    if (typeof window === 'undefined') {
        return;
    }

    localStorage.removeItem(name);
};

export function setSettings(name: string, value: string): void {
    if (typeof window === 'undefined') {
        return;
    }

    let cookieValue = getCookie('settings');
    let localValue = JSON.stringify({[name]: value});

    if (cookieValue) {
        localValue = JSON.stringify({...JSON.parse(cookieValue), [name]: JSON.parse(value)});
    }

    setLocalStorage(name, value);
    setCookie('settings', localValue, 365);
}
