import { IAssetManager } from "./assets";
import { Server } from "./server";
import { AssetManager as V1 } from "./v1";
import { AssetManager as V2 } from "./v2";
import { AssetManager as V3 } from "./v3";
import { AssetManager as V4 } from "./v4";
import { AssetManager as V5 } from "./v5";
import mysql from "mysql2/promise";

async function main() {
  const port = parseInt(process.env.SERVER_PORT || "8080", 10);
  const config = {
    host: process.env.DB_HOST || "localhost",
    port: parseInt(process.env.DB_PORT || "3306", 10),
    user: process.env.DB_USER_NAME || "root",
    password: process.env.DB_USER_PASSWORD || "rootpwd",
    database: process.env.DB_SCHEMA || "zero-downtime-training",
  };
  const pool = mysql.createPool(config);

  let assetMngr: IAssetManager;
  switch (process.env.APP_VERSION) {
    case "v1":
      assetMngr = new V1(pool);
      break;
    case "v2":
      assetMngr = new V2(pool);
      break;
    case "v3":
      assetMngr = new V3(pool);
      break;
    case "v4":
      assetMngr = new V4(pool);
      break;
    case "v5":
      assetMngr = new V5(pool);
      break;
    default:
      throw new Error(`Unsupported APP_VERSION: ${process.env.APP_VERSION}`);
  }

  const server = new Server(assetMngr, port);
  await server.start();
}

void main();
