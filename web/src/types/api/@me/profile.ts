import type { SimpleResponse } from "@/types/index";

export interface EmailUpdate extends SimpleResponse {}
export interface EmailSendVerify extends SimpleResponse {}
export interface EmailCodeVerify extends SimpleResponse {}

export interface PasswordUpdate extends SimpleResponse {}
export interface NicknameUpdate extends SimpleResponse {}
