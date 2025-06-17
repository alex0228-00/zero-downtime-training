import express, { Request, Response } from "express";

export interface IAsset {
  id: string;
  name: string;
  source: string;
}

export interface IAssetManager {
  migration(): Promise<void>;
  createAsset(asset: IAsset): Promise<IAsset>;
  readAssetByID(id: string): Promise<IAsset | null>;
  updateSourceByID(id: string, source: string): Promise<void>;
  deleteAssetByID(id: string): Promise<void>;
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
      const created = await this.assetManager.createAsset(asset);
      res.status(200).json(created);
    } catch (error) {
      res.status(500).send({ error: "Failed to create asset" });
    }
  }

  async handleReadAssetByID(req: Request, res: Response): Promise<void> {
    try {
      const id = req.params.id;
      const asset = await this.assetManager.readAssetByID(id);
      if (asset) {
        res.status(200).send(asset);
      } else {
        res.status(404).send({ error: "Asset not found" });
      }
    } catch (error) {
      res.status(500).send({ error: "Failed to read asset" });
    }
  }

  async handleUpdateSourceByID(req: Request, res: Response): Promise<void> {
    try {
      const id = req.params.id;
      const { source } = req.body;
      await this.assetManager.updateSourceByID(id, source);
      res.status(200).send({ message: "Asset source updated successfully" });
    } catch (error) {
      res.status(500).send({ error: "Failed to update asset source" });
    }
  }

  async handleDeleteAssetByID(req: Request, res: Response): Promise<void> {
    try {
      const id = req.params.id;
      await this.assetManager.deleteAssetByID(id);
      res.status(200).send({ message: "Asset deleted successfully" });
    } catch (error) {
      res.status(500).send({ error: "Failed to delete asset" });
    }
  }
}
