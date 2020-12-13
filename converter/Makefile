CC := gcc
LIBS := -lprotobuf-c -lassimp
CFLAGS := -g
PROTOCFLAGS := --experimental_allow_proto3_optional

OBJS = object.pb-c.o main.o
TARGET = convert

default: object_pb2.py
all: $(TARGET)

.PHONY: clean

%_pb2.py: %.proto
	protoc --python_out=. $(PROTOCFLAGS) $<

%.pb-c.c: %.proto
	protoc-c --c_out=. $(PROTOCFLAGS) $<

%.o: %.c
	$(CC) -c -o $@ $< $(CFLAGS)

$(TARGET): $(OBJS)
	$(CC) -o $@ $^ $(LIBS) $(CFLAGS)

clean:
	rm -f $(TARGET)
	rm -f $(OBJS)
	rm -f *.pb-c.*
	rm -f *_pb2.py
