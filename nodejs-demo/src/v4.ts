import mysql from "mysql2/promise";
import { AssetManager as V3 } from "./v3";
import { IAsset } from "./assets";

export class AssetManager extends V3 {
  constructor(pool: mysql.Pool) {
    super(pool);
  }

  async migration(): Promise<void> {}

  async createAsset(asset: IAsset): Promise<IAsset> {
    const { id, name, source } = asset;

    await this.withTransaction(async (conn) => {
      await this.insertSourceIfNotExist(conn, source);
      await conn.query(
        `
          INSERT INTO assets (id, name, source_id)
          SELECT ?, ?, id
          FROM sources
          WHERE name = ?;
        `,
        [id, name, source]
      );
    });
    return asset;
  }

  async updateSourceByID(id: string, source: string): Promise<void> {
    await this.withTransaction(async (conn) => {
      await this.insertSourceIfNotExist(conn, source);
      await conn.query(
        `
          UPDATE assets a
          JOIN (
            SELECT id, ? as asset_id 
            FROM sources 
            WHERE name = ?
          ) AS s
          ON 
            s.asset_id = a.id
          SET 
            a.source_id = s.id;
        `,
        [id, source]
      );
    });
  }
}
