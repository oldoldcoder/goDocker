#!/bin/bash

rm -rf overlay/
# 创建实验目录结构
mkdir -p overlay/{mnt,workdir,container-layer,image-layer{1..4}}

# 在各层中创建测试文件
echo "I am container layer" > overlay/container-layer/container-layer.txt
for i in {1..4}; do
    echo "I am image layer${i}" > overlay/image-layer${i}/image-layer${i}.txt
done

# 挂载 OverlayFS（需要 root 权限）
echo -e "\n=== 挂载 OverlayFS ==="
sudo mount -t overlay overlay -o lowerdir=overlay/image-layer4:overlay/image-layer3:overlay/image-layer2:overlay/image-layer1,upperdir=overlay/container-layer,workdir=overlay/workdir overlay/mnt

# 检查挂载结果
echo -e "\n=== 挂载后文件系统 ==="
df -h | grep overlay
ls -l overlay/mnt/

echo -e "\n=== 实验完成，目录结构 ==="
tree overlay
