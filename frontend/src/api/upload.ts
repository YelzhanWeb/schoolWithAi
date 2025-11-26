import api from './axios';

interface UploadResponse {
    url: string;
}

export const uploadApi = {
    uploadFile: async (file: File, type: 'avatar' | 'cover' | 'lesson'): Promise<string> => {
        const formData = new FormData();
        formData.append('file', file);
        formData.append('type', type);

        const response = await api.post<UploadResponse>('/upload', formData, {
            headers: {
                'Content-Type': 'multipart/form-data',
            },
        });
        
        return response.data.url;
    }
};