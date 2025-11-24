import type { Position } from "./Position";
import type { TileType } from "./TileType";

export interface Tile {
  tileType: TileType;
  position: Position;
}

export interface FrameProps {
  tiles: Tile[];
  title: string;
}