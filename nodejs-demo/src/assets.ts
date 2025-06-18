import mysql from "mysql2/promise";

export interface IAsset {
  id: string;
  name: string;
  source: string;
}

export abstract class IAssetManager {
  constructor(public pool: mysql.Pool) {}

  abstract migration(): Promise<void>;
  abstract createAsset(asset: IAsset): Promise<IAsset>;
  abstract readAssetByID(id: string): Promise<IAsset | null>;
  abstract updateSourceByID(id: string, source: string): Promise<void>;
  abstract deleteAssetByID(id: string): Promise<void>;

  async withTransaction(
    callBack: (conn: mysql.PoolConnection) => Promise<void>
  ): Promise<void> {
    const conn = await this.pool.getConnection();
    try {
      await conn.beginTransaction();
      await callBack(conn);
      await conn.commit();
    } catch (err) {
      await conn.rollback();
      throw err;
    } finally {
      conn.release();
    }
  }
}