import { InfoResponse } from "./responses";
import { Info } from "./types";

export const getInfo = async (): Promise<Info> => {
  const res = await fetch('/api/info')
  if (res.status !== 200) {
    const txt = await res.text()
    throw new Error(txt)
  }
  const infoRes = await res.json() as InfoResponse
  return new Info(infoRes)
}
