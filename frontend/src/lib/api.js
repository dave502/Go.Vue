import axios from 'axios';
import camelCaseKeys from 'camelcase-keys';
import snakeCaseKeys from 'snakecase-keys';



function isObject(value) {
    return typeof value === 'object' && value instanceof Object;
}

export function transformSnakeCase(data) {
    if (isObject(data) || Array.isArray(data)) {
        return snakeCaseKeys(data, { deep: true });
    }
    if (typeof data === 'string') {
        try {
            const parsedString = JSON.parse(data);
            const snakeCase = snakeCaseKeys(parsedString, { deep: true });
            return JSON.stringify(snakeCase);
        } catch (error) {
            // Bailout with no modification
            return data;
        }
    }
    return data;
}
export function transformCamelCase(data) {
    if (isObject(data) || Array.isArray(data)) {
        return camelCaseKeys(data, { deep: true });
    }
    return data;
}

// export axios object lets call axios's methods as api.<method_name>
export default axios.create({
    baseURL: import.meta.env.VITE_BACKEND_URL || "http://0.0.0.0:9010",
    // withCredentials indicates whether or not crosssite access control requests 
    // should be made using credentials such as cookies and authorization headers.
    withCredentials: true,
    transformRequest: [...axios.defaults.transformRequest,
                                            transformSnakeCase],
    transformResponse: [...axios.defaults.transformResponse,
                                            transformCamelCase],
});

