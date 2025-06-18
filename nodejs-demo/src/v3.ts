import mysql from "mysql2/promise";
import { AssetManager as V2 } from "./v2";
import { IAsset } from "./assets";

export class AssetManager extends V2 {
  constructor(pool: mysql.Pool) {
    super(pool);
  }

  async migration(): Promise<void> {
    await this.pool.query(
      `
        UPDATE assets a
        JOIN sources s 
        ON a.source = s.name
        SET 
            a.source_id = s.id;
      `
    );
  }

  async readAssetByID(id: string): Promise<IAsset | null> {
    const [rows] = await this.pool.query(
      `
        SELECT a.id, a.name, s.name AS source
        FROM assets a
        JOIN sources s ON a.source_id = s.id
        WHERE a.id = ?
      `,
      [id]
    );
    const results = rows as IAsset[];
    return results[0] ?? null;
  }
}
