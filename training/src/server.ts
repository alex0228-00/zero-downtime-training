import express, { Request, Response } from "express";
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

export class Server {
  constructor(private assetManager: IAssetManager) {}

  async start() {
    await this.assetManager.migration();

    const app = express();
    this.registerHandlers(app);

    const port = process.env.SERVER_PORT || 80;
    app.listen(port, () => {
      console.log(`Server is running at http://localhost:${port}`);
    });
  }

  registerHandlers(app: express.Express) {
    app.use(express.json());

    // Health check
    app.get("/health", (_req: Request, res: Response) => {
      res.status(200).send("OK");
    });

    app.post("/api/asset", this.handleCreateAsset.bind(this));
    app.put("/api/asset/:id", this.handleUpdateSourceByID.bind(this));
    app.get("/api/asset/:id", this.handleReadAssetByID.bind(this));
    app.delete("/api/asset/:id", this.handleDeleteAssetByID.bind(this));
  }

  async handleCreateAsset(req: Request, res: Response): Promise<void> {
    try {
      const asset: IAsset = req.body;

      console.log(`Creating asset: ${JSON.stringify(asset)}`);
      const created = await this.assetManager.createAsset(asset);
      res.status(200).json(created);
    } catch (error) {
      logAndReturnError(res, "Failed to create asset", error);
    }
  }

  async handleReadAssetByID(req: Request, res: Response): Promise<void> {
    try {
      const id = req.params.id;

      console.log(`Reading asset from db, id=${id}`);
      const asset = await this.assetManager.readAssetByID(id);

      console.log(`Asset read: ${JSON.stringify(asset)}`);
      if (asset) {
        res.status(200).send(asset);
      } else {
        res.status(404).send({ error: "Asset not found" });
      }
    } catch (error) {
      logAndReturnError(res, "Failed to read asset", error);
    }
  }

  async handleUpdateSourceByID(req: Request, res: Response): Promise<void> {
    try {
      const id = req.params.id;
      const { source } = req.body;

      console.log(`Updating asset source, id=${id}, source=${source}`);
      await this.assetManager.updateSourceByID(id, source);
      res.status(200).send({ message: "Asset source updated successfully" });
    } catch (error) {
      logAndReturnError(res, "Failed to update asset", error);
    }
  }

  async handleDeleteAssetByID(req: Request, res: Response): Promise<void> {
    try {
      const id = req.params.id;
      await this.assetManager.deleteAssetByID(id);
      res.status(200).send({ message: "Asset deleted successfully" });
    } catch (error) {
      logAndReturnError(res, "Failed to delete asset", error);
    }
  }
}

function logAndReturnError(
  res: Response,
  message: string,
  error: unknown
): void {
  const errorMessage = error instanceof Error ? error.message : "Unknown error";
  console.error(errorMessage);
  res.status(500).send({ error: `${message}: ${errorMessage}` });
}
