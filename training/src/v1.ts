import { IAsset, IAssetManager } from "./server";
import mysql from "mysql2/promise";

export class AssetManager extends IAssetManager {
  constructor(pool: mysql.Pool) {
    super(pool);
  }

  async migration(): Promise<void> {
    await this.pool.query(
      `
        CREATE TABLE IF NOT EXISTS assets (
          id VARCHAR(36) PRIMARY KEY,
          name VARCHAR(255),
          source VARCHAR(255)
        )
      `
    );
  }

  async createAsset(asset: IAsset): Promise<IAsset> {
    const { id, name, source } = asset;
    await this.pool.query(
      `INSERT INTO assets (id, name, source) VALUES (?, ?, ?)`,
      [id, name, source]
    );
    return asset;
  }

  async readAssetByID(id: string): Promise<IAsset | null> {
    const [rows] = await this.pool.query(
      `
        SELECT * 
        FROM assets 
        WHERE id = ?
      `,
      [id]
    );
    const results = rows as IAsset[];
    return results[0] ?? null;
  }

  async updateSourceByID(id: string, source: string): Promise<void> {
    await this.pool.query(
      `
        UPDATE assets 
        SET source = ? 
        WHERE id = ?
      `,
      [source, id]
    );
  }

  async deleteAssetByID(id: string): Promise<void> {
    await this.pool.query("DELETE FROM assets WHERE id = ?", [id]);
  }
}
