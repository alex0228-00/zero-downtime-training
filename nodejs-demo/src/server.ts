import express, { Request, Response } from "express";

import { IAsset, IAssetManager } from "./assets";

export class Server {
  constructor(private assetManager: IAssetManager, private port: number) {}

  async start() {
    await this.assetManager.migration();

    const app = express();
    this.registerHandlers(app);

    app.listen(this.port, "0.0.0.0", () => {
      console.log(`Server is running at http://localhost:${this.port}`);
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
