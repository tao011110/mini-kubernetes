SOURCE = job.cu

TARGETBIN := ./cuda

CC = nvcc

$(TARGETBIN):$(SOURCE)
	$(CC)  $(SOURCE) -o $(TARGETBIN)

.PHONY:clean
clean:
	-rm -rf $(TARGETBIN)
