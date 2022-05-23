#include <stdio.h>
#include "cuda_runtime.h"
#include "device_launch_parameters.h"
#include <cuda_runtime.h>
#include <stdio.h>
#include <iostream>

using namespace std;

#define CHECK(call)                                                    \
  {                                                                    \
    const cudaError_t error = call;                                    \
    if (error != cudaSuccess) {                                        \
      printf("ERROR: %s:%d,", __FILE__, __LINE__);                     \
      printf("code:%d,reason:%s\n", error, cudaGetErrorString(error)); \
      exit(1);                                                         \
    }                                                                  \
  }

void initDevice(int devNum) {
  int dev = devNum;
  cudaDeviceProp deviceProp;
  CHECK(cudaGetDeviceProperties(&deviceProp, dev));
  printf("Using device %d: %s\n", dev, deviceProp.name);
  CHECK(cudaSetDevice(dev));
}

void initialData(float* ip, int size) {
  time_t t;
  srand((unsigned)time(&t));
  for (int i = 0; i < size; i++)
    ip[i] = (float)(rand() & 0xffff) / 1000.0f;
}

__global__ void sumMatrix(float* MatA, float* MatB, float* MatC, int nx,
                          int ny) {
  int ix = threadIdx.x + blockDim.x * blockIdx.x;
  int iy = threadIdx.y + blockDim.y * blockIdx.y;
  int idx = ix + iy * ny;
  if (ix < nx && iy < ny)
    MatC[idx] = MatA[idx] + MatB[idx];
}

int main(int argc, char** argv) {
  //init dev
  initDevice(0);

  int nx = 1 << 12;
  int ny = 1 << 12;
  int nBytes = nx * ny * sizeof(float);

  float* A_host = (float*)malloc(nBytes);
  float* B_host = (float*)malloc(nBytes);
  float* C_from_gpu = (float*)malloc(nBytes);
  initialData(A_host, nx * ny);
  initialData(B_host, nx * ny);

  float* A_dev = NULL;
  float* B_dev = NULL;
  float* C_dev = NULL;
  CHECK(cudaMalloc((void**)&A_dev, nBytes));
  CHECK(cudaMalloc((void**)&B_dev, nBytes));
  CHECK(cudaMalloc((void**)&C_dev, nBytes));

  CHECK(cudaMemcpy(A_dev, A_host, nBytes, cudaMemcpyHostToDevice));
  CHECK(cudaMemcpy(B_dev, B_host, nBytes, cudaMemcpyHostToDevice));

  dim3 threadsPerBlock(32, 32);
  cout << "threadsPerBlock.x = " << threadsPerBlock.x << endl;
  cout << "threadsPerBlock.y = " << threadsPerBlock.y << endl;

  dim3 numBlocks((nx - 1) / threadsPerBlock.x + 1,
                 (ny - 1) / threadsPerBlock.y + 1);
  cout << "numBlocks.x = " << numBlocks.x << "   numBlocks.y=" << numBlocks.y
       << endl;

  sumMatrix<<<numBlocks, threadsPerBlock>>>(A_dev, B_dev, C_dev, nx, ny);
  CHECK(cudaDeviceSynchronize());

  CHECK(cudaMemcpy(C_from_gpu, C_dev, nBytes, cudaMemcpyDeviceToHost));

  cudaFree(A_dev);
  cudaFree(B_dev);
  cudaFree(C_dev);
  free(A_host);
  free(B_host);
  free(C_from_gpu);
  cudaDeviceReset();

  return 0;
}
