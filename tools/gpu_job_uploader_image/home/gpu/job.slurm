#!/bin/bash

#SBATCH --job-name=GPU_JOB
#SBATCH --partition=a100
#SBATCH -N 1
#SBATCH --cpus-per-task=16
#SBATCH --ntasks-per-node=1
#SBATCH --gres=gpu:1
#SBATCH --output=result.out
#SBATCH --error=error.err
#SBATCH --time=00-00:2:00

module load gcc cuda

make
./cuda
