import axios from 'axios';
import authHeader from './auth-header';

const API_URL = process.env.VUE_APP_HTTP_SERVER_URL;

class FeedService {
    addFeed(feed) {
        return axios.post(API_URL + 'feed/add', { body: feed }, { headers: authHeader() });
    }

    getFeed(userId) {
        return axios.get(API_URL + 'feed/get/' + userId);
    }
}

export default new FeedService();
