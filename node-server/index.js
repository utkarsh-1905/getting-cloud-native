require("dotenv").config();
const express = require("express");
const app = express();
const cors = require("cors");
const morgan = require("morgan");
const responseTime = require("response-time");
const mongoose = require("mongoose");
const port = process.env.PORT || 3000;

app.use(cors());
app.use(morgan("dev"));
app.use(responseTime());
app.use(require("express-status-monitor")());

mongoose.connect(process.env.MONGODB, {}, (err) => {
  if (!err) {
    console.log("Connected to MongoDB");
  } else {
    console.log(err);
  }
});

const randomSchema = new mongoose.Schema({
  age: Number,
  message: String,
});

const TestNode = mongoose.model("TestNode", randomSchema);

app.get("/", (req, res) => {
  res.status(200).json({
    status: "OK",
  });
});

app.get("/db/create", async (req, res) => {
  try {
    for (let i = 0; i < 100; i++) {
      const test = new TestNode({
        age: Math.floor(Math.random() * 100),
        message: Math.random().toString(),
      });
      await test.save();
    }
    res.status(200).json({ status: "OK" });
  } catch (e) {
    res.status(500).json({ status: "ERROR" });
  }
});

app.get("/db/read", async (req, res) => {
  try {
    const test = await TestNode.find({ age: req.query.age });
    res.status(200).json(test);
  } catch (e) {
    res.status(500).json({ status: "ERROR" });
  }
});

app.get("/db/delete", async (req, res) => {
  try {
    await TestNode.deleteMany({});
    res.status(200).json({ status: "OK" });
  } catch (e) {
    res.status(500).json({ status: "ERROR" });
  }
});

app.listen(port, () => console.log(`Server is running on port ${port}`));
