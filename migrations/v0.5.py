import os 

from pymongo import MongoClient

MONGO_URI = os.getenv("STUFFAI_API_MONGO_URI")
client = MongoClient(MONGO_URI)
db = client.stuffai

imgs = db.images.find()
for img in imgs:
    db.jobs.update_one({"_id": img["_id"]}, {"$set": {"listeners": [img["user"]["_id"]]}})