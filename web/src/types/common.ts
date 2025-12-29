export interface SqlNullInt64 {
    Int64: number;
    Valid: boolean;
}

export interface SqlNullInt32 {
    Int32: number;
    Valid: boolean;
}

export interface SqlNullFloat64 {
    Float64: number;
    Valid: boolean;
}

export interface SqlNullString {
    String: string;
    Valid: boolean;
}

export interface SqlNullTime {
    Time: string;
    Valid: boolean;
}

export function unwrap<T>(val: { Valid: boolean, [key: string]: any } | undefined | null, key: string, fallback: T): T {
    if (!val || !val.Valid) return fallback;
    return val[key] as T;
}
