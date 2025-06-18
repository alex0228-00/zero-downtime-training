import { IAssetManager } from "./assets";
import { Server } from "./server";
import { AssetManager as V1 } from "./v1";
import { AssetManager as V2 } from "./v2";
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
    default:
      throw new Error(`Unsupported APP_VERSION: ${process.env.APP_VERSION}`);
  }

  const server = new Server(assetMngr, port);
  await server.start();
}

void main();
