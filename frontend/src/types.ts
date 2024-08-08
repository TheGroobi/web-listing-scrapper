export type ApiResult<T = undefined> = {
    statusCode: number;
    ok: boolean;
    data: T;
} & ApiMessage

type ApiMessage = ({ message: string; error?: never; } | { error: string; message?: never; });

export type OtomotoData = {
    Title: string
    Link: string
    Gearbox: string
    BodyType: string
    FuelType: string
    Color: string
    Version: string | undefined
    Year: number
    Power: number
    Mileage: number
    ID: number
    Price: number
}
