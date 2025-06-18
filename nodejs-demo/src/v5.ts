import mysql from "mysql2/promise";
import { AssetManager as V4 } from "./v4";

export class AssetManager extends V4 {
  constructor(pool: mysql.Pool) {
    super(pool);
  }

  async migration(): Promise<void> {
    await this.pool.query("ALTER TABLE assets DROP COLUMN source;");
  }
}
