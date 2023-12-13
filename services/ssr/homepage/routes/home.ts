import express from 'express';
import * as homeControllers from '../controllers/home';
const router = express.Router();

router.get('/', homeControllers.renderHomepage);

export default router;
