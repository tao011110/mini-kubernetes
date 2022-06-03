# cuda矩阵乘法/矩阵加法代码解析

## 程序流程
- 使用`cudaGetDeviceProperties`和`cudaSetDevice`初始化和检查GPU设备
- 在CPU内存中初始化矩阵A, B
- 使用`cudaMalloc`在显存中malloc两个输入矩阵的空间, 并使用`cudaMemcpy`将A, B拷贝进显存对应位置
- (并行加速的核心)定义线程块和线程块内部线程模型, 调用核函数计算(考虑到矩阵的特点, 使用了二维线程模型)
````c
// 矩阵乘法
dim3 threadPerBlock(16, 16);
dim3 blockNumber((Col+threadPerBlock.x-1)/ threadPerBlock.x, (Row+threadPerBlock.y-1)/ threadPerBlock.y );
matrix_mul_gpu <<<blockNumber, threadPerBlock >>> (d_dataA, d_dataB, d_dataC, Col);

dim3 threadPerBlock(32, 32);
dim3 numBlocks((nx - 1) / threadsPerBlock.x + 1, (ny - 1) / threadsPerBlock.y + 1);
sumMatrix<<<numBlocks, threadPerBlock>>>(A_dev, B_dev, C_dev, nx, ny);
````
- 使用`cudaMemcpy`将计算结果从显存中拷贝进主存中
- 使用`free`释放主存空间, 使用`cudaFree`释放显存空间

## 核函数
- 矩阵加法(每个线程处理矩阵上一个坐标位置的加法), 可以证明, 一个线程可以唯一确定一个(x,y), 且可以覆盖矩阵内的全部(x, y)组合
````c
__global__ void sumMatrix(float* MatA, float* MatB, float* MatC, int nx, int ny) {
  int y = threadIdx.x + blockDim.x * blockIdx.x; //可以被认为是列号
  int x = threadIdx.y + blockDim.y * blockIdx.y; //可以被认为是行号
  int idx = y + x * nx; //矩阵用一维数组形式存储, iy行ix列对应索引为iy*ny+ix处的位置
  if (y < ny && x < nx)
    MatC[idx] = MatA[idx] + MatB[idx];
}
````
- 矩阵(方阵)乘法(每个线程处理结果上一个位置的值, 即第i行和第j列的点乘)
````c
__global__ void matrix_mul_gpu(int *M, int* N, int* P, int width){
    int i = threadIdx.x + blockDim.x * blockIdx.x;
    int j = threadIdx.y + blockDim.y * blockIdx.y;
    int sum = 0;
    for(int k = 0; k < width; k++)
        sum += M[j*width+k] * N[k*width+i];
    P[j*width+i] = sum;
}
````
