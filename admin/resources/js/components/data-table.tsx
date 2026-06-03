"use client"

import {
    ColumnDef,
    flexRender,
    getCoreRowModel,
    VisibilityState,
    useReactTable,
} from "@tanstack/react-table";
import {
    DropdownMenu,
    DropdownMenuCheckboxItem,
    DropdownMenuItem,
    DropdownMenuContent,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import {
    ChevronLeft,
    ChevronRight,
    ChevronsLeft,
    ChevronsRight,
    CirclePlus,
    Columns3,
    Search,
    X,
    Eye,
    EyeOff,
    Check,
    ChevronDown,
    Settings2
} from "lucide-react";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import { Separator } from "@/components/ui/separator"
import { Badge } from "@/components/ui/badge";
import { useState, useCallback, useEffect } from "react";
import { Label } from '@/components/ui/label';
import { Button } from "@/components/ui/button";
import { router, usePage } from "@inertiajs/react";
import { Input } from '@/components/ui/input';
import {setSettings} from "@/hooks/use-settings";

interface DataTableProps<TData, TValue> {
    columns: ColumnDef<TData, TValue>[]
    typeKey: string
    filters?: ColumnDef<TData, TValue>[]
}

export function DataTable<TData, TValue>({columns, typeKey, filters}: DataTableProps<TData, TValue>) {
    const { props: { data = [], search = "", activeFilters = {}, dataFilters = [], auth }, url = "" } = usePage();
    const pagination = {
        show: data.hasOwnProperty('data'),
        current_page: data.hasOwnProperty('current_page') ? data.current_page : 0,
        last_page: data.hasOwnProperty('last_page') ? data.last_page : 0,
        per_page: data.hasOwnProperty('per_page') ? data.per_page : 0,
    };
    const IS_REMEMBER_SETTINGS = auth.user.settings.is_remember_settings ?? false;
    const KEY_VIEW_COLUMNS = `settings_${typeKey}_columns`
    const KEY_PER_PAGE = `settings_${typeKey}_per_page`

    const dataTable = Array.isArray(data) ? data : (pagination.show ? data.data : []);
    const route = url.split('?')[0]

    const [columnVisibility, setColumnVisibility] = useState<{}>(JSON.parse(localStorage.getItem(KEY_VIEW_COLUMNS) ?? '{}'))
    const [searchValue, setSearchValue] = useState(search)
    const [isSearching, setIsSearching] = useState(false)
    const [composition, setComposition] = useState(activeFilters)

    const table = useReactTable({
        data: dataTable,
        columns: columns,
        getCoreRowModel: getCoreRowModel(),
        onColumnVisibilityChange: ( updater ) => {
            setColumnVisibility((previousState) => {
                const newState = (typeof updater === 'function') ? updater(previousState) : updater;
                if (IS_REMEMBER_SETTINGS) {
                    const value = JSON.stringify(newState);
                    setSettings(KEY_VIEW_COLUMNS, value)
                }

                return newState
            })
        },
        state: {
            columnVisibility: columnVisibility,
        },
        manualFiltering: true,
        manualPagination: pagination.show,
        pageCount: pagination.last_page,
        initialState: {
            pagination: {
                pageIndex: pagination.current_page - 1,
                pageSize: pagination.per_page,
            }
        }
    })

    const handlePageChange = (page: number) => {
        const params = new URLSearchParams();
        params.set('page', page.toString());
        params.set('perPage', pagination.per_page.toString());
        if (searchValue) {
            params.set('search', searchValue);
        }

        router.get(`${route}?${params.toString()}`)
    }

    const handlePerPageChange = (perPage: number) => {
        const params = new URLSearchParams();
        params.set('page', '1');
        params.set('perPage', perPage.toString());
        if (searchValue) {
            params.set('search', searchValue);
        }

        if (IS_REMEMBER_SETTINGS) {
            const value = JSON.stringify(perPage);
            setSettings(KEY_PER_PAGE, value)
        }

        router.get(`${route}?${params.toString()}`, {}, {
            preserveState: true,
            replace: true
        })
    }

    const handleSearch = (value: string) => {
        setSearchValue(value)
        setIsSearching(true)
        const params = new URLSearchParams();
        params.set('page', '1');
        params.set('perPage', pagination.per_page.toString());
        if (value.trim()) {
            params.set('search', value.trim());
        }
        params.set('filters', JSON.stringify(composition));
        router.get(`${route}?${params.toString()}`, {}, {
            onFinish: () => setIsSearching(false),
            preserveState: true,
            replace: true
        })
    }

    const handleFilter = (e, key, item) => {
        e.preventDefault();
        let result;
        setComposition(prev => {
            const currentArray = prev[key] ?? [];
            if (currentArray.includes(item)) {
                result = {
                    ...prev,
                    [key]: currentArray.filter(el => el !== item)
                };
            } else {
                result = {
                    ...prev,
                    [key]: [...currentArray, item]
                };
            }
            return result;
        });
        const params = new URLSearchParams();
        params.set('page', '1');
        params.set('perPage', pagination.per_page.toString());
        if (searchValue) {
            params.set('search', searchValue);
        }
        params.set('filters', JSON.stringify(result));
        router.get(`${route}?${params.toString()}`, {}, {
            preserveState: true,
            replace: true
        })
    }

    const handleReset = () => {
        setSearchValue('')
        setComposition({})
        const params = new URLSearchParams();
        params.set('page', '1');
        params.set('perPage', pagination.per_page.toString());
        params.delete('search');
        params.delete('filters');

        router.get(`${route}?${params.toString()}`, {}, {
            preserveState: true,
        })
    }

    return (
        <div>
            <div className="grid lg:flex items-center py-4 gap-2">
                <div className="flex items-center gap-2">
                    <div className="relative flex-1 max-w-sm">
                        <Search className="absolute left-2 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                        <Input
                            placeholder="Пошук..."
                            value={searchValue}
                            onChange={(event) => setSearchValue(event.target.value)}
                            onKeyDown={(event) => {
                                if (event.key === 'Enter') {
                                    handleSearch(searchValue)
                                }
                            }}
                            disabled={isSearching}
                            className="pl-8 pr-10 py-0 h-8 w-full"
                        />
                        {searchValue && (
                            <button
                                onClick={() => handleSearch('')}
                                className="absolute right-2 top-1/2 -translate-y-1/2 text-muted-foreground text-sm hover:text-foreground items-center justify-center"
                                disabled={isSearching}
                            >
                                <X className="h-4 w-4" />
                            </button>
                        )}
                    </div>
                    <Button
                        variant="secondary"
                        size="sm"
                        onClick={() => handleSearch(searchValue)}
                        className="font-normal text-sm"
                        disabled={isSearching}
                    >
                        Пошук
                    </Button>
                </div>
                { dataFilters && <div className="flex items-center gap-2">
                    { filters?.map((filter) => {
                        return (
                            <DropdownMenu key={filter.accessorKey}>
                                <DropdownMenuTrigger asChild className="hidden sm:flex">
                                    <Button variant="outline" size="sm" className="gap-1 text-sm flex items-center font-normal">
                                        <CirclePlus className="h-4 w-4" />
                                        {filter.header}
                                        {( composition[filter.accessorKey] ?? [] ).length > 0 && <Separator orientation="vertical" />}
                                        {
                                            ( composition[filter.accessorKey] ?? [] ).length > 2
                                                ? <Badge size="xs" variant="secondary" className="font-normal">Обрано {composition[filter.accessorKey].length}</Badge>
                                                : <>{( composition[filter.accessorKey] ?? [] ).map((element, index) => <Badge key={index} size="xs" variant="secondary" className="font-normal">{element}</Badge>)}</>
                                        }
                                    </Button>
                                </DropdownMenuTrigger>
                                <DropdownMenuContent>
                                    { dataFilters[filter.accessorKey].map((element, index) =>
                                        <DropdownMenuItem
                                            key={index}
                                            onClick={(e) => handleFilter(e, filter.accessorKey, element)}
                                        >
                                            {(composition[filter.accessorKey] ?? []).includes(element) ? (
                                                <div className="flex mr-4 size-4 items-center justify-center rounded-[4px] border border-primary bg-primary text-primary-foreground">
                                                    <Check className="lucide lucide-check size-3.5 text-primary-foreground"/>
                                                </div>
                                            ) : (
                                                <div className="flex mr-4 size-4 items-center justify-center rounded-[4px] border border-input [&_svg]:invisible" />
                                            )}
                                            { typeof filter.item === 'function' ? filter.item(element) : element }
                                        </DropdownMenuItem>)
                                    }
                                </DropdownMenuContent>
                            </DropdownMenu>)
                        })
                    }
                </div> }
                {
                    ( searchValue || filters?.filter((filter) => (composition[filter.accessorKey] ?? []).length > 0).length > 0 )
                    &&
                        <Button
                            variant="secondary"
                            size="sm"
                            className="gap-1 text-sm font-normal"
                            onClick={ handleReset }
                        >
                            Скинути <X className="h-4 w-4" />
                        </Button>
                }

                <DropdownMenu>
                    <DropdownMenuTrigger asChild className="hidden sm:flex">
                        <Button variant="outline" className="ml-auto">
                            <Settings2 className="h-4 w-4" />
                            Колонки
                            { Object.values(columnVisibility).includes(false) && <EyeOff className="h-2 w-2 text-red-500" /> }
                            <ChevronDown />
                        </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                        {table
                            .getAllColumns()
                            .filter(
                                (column) => column.getCanHide()
                            )
                            .map((column) => {
                                let checked = column.getIsVisible()
                                return (
                                    <DropdownMenuCheckboxItem
                                        key={column.id}
                                        className="capitalize"
                                        checked={checked}
                                        onCheckedChange={(value) =>
                                            column.toggleVisibility(value)
                                        }
                                    >
                                        {checked ? <Eye className="h-2 w-2" /> : <EyeOff className="h-2 w-2 text-red-500" />}
                                        {column.id}
                                    </DropdownMenuCheckboxItem>
                                )
                            })}
                    </DropdownMenuContent>
                </DropdownMenu>
            </div>
            <div className="overflow-hidden rounded-md border">
                <Table>
                    <TableHeader>
                        {table.getHeaderGroups().map((headerGroup) => (
                            <TableRow key={headerGroup.id} className="bg-muted/70">
                                {headerGroup.headers.map((header) => {
                                    return (
                                        <TableHead key={header.id}>
                                            {header.isPlaceholder
                                                ? null
                                                : flexRender(
                                                    header.column.columnDef.header,
                                                    header.getContext()
                                                )}
                                        </TableHead>
                                    )
                                })}
                            </TableRow>
                        ))}
                    </TableHeader>
                    <TableBody>
                        {table.getRowModel().rows?.length ? (
                            table.getRowModel().rows.map((row) => (
                                <TableRow key={row.id}>
                                    {row.getVisibleCells().map((cell) => (
                                        <TableCell key={cell.id} className="whitespace-normal">
                                            {flexRender(cell.column.columnDef.cell, cell.getContext())}
                                        </TableCell>
                                    ))}
                                </TableRow>
                            ))
                        ) : (
                            <TableRow>
                                <TableCell colSpan={columns.length} className="h-24 text-center">
                                    {searchValue ? "Збігів не знайдено" : "Дані відсутні"}.
                                </TableCell>
                            </TableRow>
                        )}
                    </TableBody>
                </Table>
            </div>
            { pagination.show && <div className="flex items-center justify-end space-x-2 py-4">
                <div className="flex w-full items-center gap-8 lg:w-fit">
                    <div className="hidden items-center gap-2 lg:flex">
                        <Label htmlFor="rows-per-page" className="text-sm font-medium">
                            Рядків на сторінці
                        </Label>
                        <Select
                            value={`${pagination.per_page}`}
                            onValueChange={(value) => handlePerPageChange(Number(value))}
                        >
                            <SelectTrigger size="sm" className="w-20" id="rows-per-page">
                                <SelectValue placeholder={pagination.per_page}/>
                            </SelectTrigger>
                            <SelectContent side="top">
                                {[10, 25, 50, 100].map((pageSize) => (
                                    <SelectItem key={pageSize} value={`${pageSize}`}>
                                        {pageSize}
                                    </SelectItem>
                                ))}
                            </SelectContent>
                        </Select>
                    </div>
                    <div className="flex w-fit items-center justify-center text-sm font-medium">
                        Сторінка {pagination.current_page} з {pagination.last_page}
                    </div>
                    <div className="ml-auto flex items-center gap-2 lg:ml-0">
                        <Button
                            variant="outline"
                            className="hidden h-8 w-8 p-0 lg:flex"
                            onClick={() => handlePageChange(1)}
                            disabled={pagination.current_page === 1}
                        >
                            <span className="sr-only">Перша сторінка</span>
                            <ChevronsLeft />
                        </Button>
                        <Button
                            variant="outline"
                            className="size-8"
                            size="icon"
                            onClick={() => handlePageChange(pagination.current_page - 1)}
                            disabled={pagination.current_page === 1}
                        >
                            <span className="sr-only">Попередня сторінка</span>
                            <ChevronLeft />
                        </Button>
                        <Button
                            variant="outline"
                            className="size-8"
                            size="icon"
                            onClick={() => handlePageChange(pagination.current_page + 1)}
                            disabled={pagination.current_page === pagination.last_page}
                        >
                            <span className="sr-only">Наступна сторінка</span>
                            <ChevronRight />
                        </Button>
                        <Button
                            variant="outline"
                            className="hidden size-8 lg:flex"
                            size="icon"
                            onClick={() => handlePageChange(pagination.last_page)}
                            disabled={pagination.current_page === pagination.last_page}
                        >
                            <span className="sr-only">Остання сторінка</span>
                            <ChevronsRight />
                        </Button>
                    </div>
                </div>
            </div>}
        </div>
    )
}
