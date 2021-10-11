function getDataFromStream(_data): Promise<any> {
    return Axios({
        url: `http://the.api.url/endpoint`,
        method: 'GET',
        onDownloadProgress: progressEvent => {
           const dataChunk = progressEvent.currentTarget.response;
           // dataChunk contains the data that have been obtained so far (the whole data so far).. 
           // So here we do whatever we want with this partial data.. 
        }
    }).then(({ data }) => Promise.resolve(data));  
}