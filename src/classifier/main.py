# Load libraries
import torch
import glob
import torch.nn as nn
from torchvision.transforms import transforms
from torch.utils.data import DataLoader
import torch.optim
from torch.autograd import Variable
import torchvision
import pathlib

from src.ConvNet import ConvNet

NUM_EPOCHS = 100

# checking for device
device = torch.device('cuda' if torch.cuda.is_available() else 'cpu')

# Transforms
transformer = transforms.Compose([
    transforms.Resize((150, 150)),
    transforms.ColorJitter(),
    transforms.RandomHorizontalFlip(),
    transforms.RandomVerticalFlip(),
    transforms.RandomRotation(+90),
    transforms.ToTensor(),  # 0-255 to 0-1, numpy to tensors
    transforms.Normalize([0.5, 0.5, 0.5],  # 0-1 to [-1,1] , formula (x-mean)/std
                         [0.5, 0.5, 0.5])
])

# Dataloader

# Path for training and testing directory
train_path = './references/train'
test_path = './references/test'

batchSize = 4

train_loader = DataLoader(
    torchvision.datasets.ImageFolder(train_path, transform=transformer),
    batch_size=batchSize, shuffle=True
)
test_loader = DataLoader(
    torchvision.datasets.ImageFolder(test_path, transform=transformer),
    batch_size=batchSize // 2, shuffle=True
)

# categories
root = pathlib.Path(train_path)
classes = sorted([j.name.split('/')[-1] for j in root.iterdir()])

# CNN Network
print(classes)

model = ConvNet(num_classes=len(classes), batchsize=batchSize).to(device)

# Optmizer and loss function
optimizer = torch.optim.Adam(model.parameters(), lr=0.001)
loss_function = nn.CrossEntropyLoss()

# calculating the size of training and testing images
train_count = len(glob.glob(train_path + '/**/*.jpg'))
test_count = len(glob.glob(test_path + '/**/*.jpg'))

print(train_count, test_count)

# Model training and saving best model

best_accuracy = 0.0

for epoch in range(NUM_EPOCHS):

    # Evaluation and training on training dataset
    model.train()
    train_accuracy = 0.0
    train_loss = 0.0

    for i, (images, labels) in enumerate(train_loader):
        if torch.cuda.is_available():
            images = Variable(images.cuda())
            labels = Variable(labels.cuda())

        optimizer.zero_grad()

        outputs = model(images)
        loss = loss_function(outputs, labels)
        loss.backward()
        optimizer.step()

        train_loss += loss.cpu().data * images.size(0)
        _, prediction = torch.max(outputs.data, 1)

        train_accuracy += int(torch.sum(prediction == labels.data))

    train_accuracy = train_accuracy / train_count
    train_loss = train_loss / train_count

    # Evaluation on testing dataset
    model.eval()

    test_accuracy = 0.0
    for i, (images, labels) in enumerate(test_loader):
        if torch.cuda.is_available():
            images = Variable(images.cuda())
            labels = Variable(labels.cuda())

        outputs = model(images)
        _, prediction = torch.max(outputs.data, 1)
        test_accuracy += int(torch.sum(prediction == labels.data))

    test_accuracy = test_accuracy / test_count

    print('Epoch: ' + str(epoch) + ' Train Loss: ' + str(train_loss) +
          ' Train Accuracy: ' + str(train_accuracy * 100) + ' Test Accuracy: ' + str(test_accuracy * 100))

    # Save the best model
    if test_accuracy > best_accuracy:
        torch.save(model.state_dict(), 'models/best_95_checkpoint_15_04_23.model')
        best_accuracy = test_accuracy

print("DONE WITH PREDICTION")
print(f"==================Best testAcc [{best_accuracy * 100} %]==================")

# checkpoint = torch.load('models/best_95_checkpoint_15_04_23.model')
# model = ConvNet(num_classes=len(classes), batchsize=batchSize)
# model.load_state_dict(checkpoint)
# model.eval()

# transformer_pred = transforms.Compose([
#     transforms.Resize((150, 150)),
#     transforms.ToTensor(),  # 0-255 to 0-1, numpy to tensors
#     transforms.Normalize([0.5, 0.5, 0.5],  # 0-1 to [-1,1] , formula (x-mean)/std
#                          [0.5, 0.5, 0.5])
# ])
# pred_path = './references/pred/'


# # prediction function

# def prediction(img_path, transformer):
#     image = Image.open(img_path)
#     image_tensor = transformer(image).float()
#     image_tensor = image_tensor.unsqueeze_(0)
#     if torch.cuda.is_available():
#         image_tensor.cuda()
#     input = Variable(image_tensor)
#     output = model(input)
#     # print(f"Output Array {output.data.numpy()}")
#     index = output.data.numpy().argmax()
#     predClass = classes[index]

#     # print(f"Input Img: {img_path[img_path.rfind('/') + 1:]}\nPredicted class {predClass}\nindex {index}",
#     #       end="\n=========\n")
#     return predClass


# images_path = glob.glob(pred_path + '/*.jpg')

# pred_accuracy = 0.0

# for i in images_path:
#     predClass = prediction(i, transformer_pred)
#     imgName = i[i.rfind('/') + 1:]
#     orgClass = imgName.split('_')[:1][0].title() + " Soil"

#     if orgClass == predClass:
#         pred_accuracy += 1

# print(f"Correct Prediction: {pred_accuracy} out of {len(images_path)} -> {pred_accuracy/len(images_path)}")
