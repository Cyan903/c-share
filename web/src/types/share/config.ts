export interface SXCU {
    Version: string;
    Name: string;
    DestinationType: string;
    RequestMethod: string;
    RequestURL: string;
    Parameters: Parameters;
    Body: string;
    FileFormName: string;
    URL: string;
    ErrorMessage: string;
}

export interface Parameters {
    token: string;
    perm: string;
    password?: string;
    comment: string;
}
