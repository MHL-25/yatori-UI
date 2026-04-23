declare module '../wailsjs/go/main/App' {
  export function GetAccountList(): Promise<string>
  export function AddAccount(jsonReq: string): Promise<string>
  export function DeleteAccount(uid: string): Promise<string>
  export function UpdateAccount(jsonReq: string): Promise<string>
  export function LoginCheck(uid: string): Promise<string>
  export function StartBrush(uid: string): Promise<string>
  export function StopBrush(uid: string): Promise<string>
  export function GetAllProgress(): Promise<string>
  export function GetProgress(uid: string): Promise<string>
  export function GetConfig(): Promise<string>
  export function SaveConfig(jsonCfg: string): Promise<string>
  export function GetPlatforms(): Promise<string>
  export function GetAiTypes(): Promise<string>
  export function MinimizeWindow(): Promise<void>
  export function ToggleMaximizeWindow(): Promise<void>
  export function CloseWindow(): Promise<void>
}

declare module '../wailsjs/runtime/runtime' {
  export function EventsOn(eventName: string, callback: (...args: any[]) => void): void
  export function EventsEmit(eventName: string, ...data: any[]): void
  export function EventsOff(eventName: string): void
}
