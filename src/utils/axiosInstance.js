import axios from 'axios';

const baseApiUrl = process.env.REACT_APP_API_BASE_URL

const axiosInstance = axios.create({
  baseURL: baseApiUrl, 
  headers: {
    'Authorization': `Bearer ${localStorage.getItem('authToken')}`,
  },
});

export default axiosInstance;
