import axios from "axios";

const baseURL = "localhost:3000";

export default axios.create({
  baseURL,
});
