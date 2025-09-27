import { useState } from 'react'
import { useForm } from 'react-hook-form'
import Modal from '../ui/Modal'
import Input from '../ui/Input'
import Button from '../ui/Button'

const AddSongModal = ({ onClose, onAddSong }) => {
  const { register, handleSubmit, formState: { errors } } = useForm()
  const [isSubmitting, setIsSubmitting] = useState(false)
  
  const onSubmit = async (data) => {
    setIsSubmitting(true)
    try {
      await onAddSong(data)
    } finally {
      setIsSubmitting(false)
    }
  }
  
  return (
    <Modal
      isOpen={true}
      onClose={onClose}
      title="Add Song to Queue"
      size="md"
    >
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        <Input
          label="Song Title"
          {...register('songName', { required: 'Song title is required' })}
          error={errors.songName?.message}
          placeholder="Enter song title"
        />
        
        <Input
          label="Artist"
          {...register('artistName', { required: 'Artist name is required' })}
          error={errors.artistName?.message}
          placeholder="Enter artist name"
        />
        
        <Input
          label="Album (Optional)"
          {...register('albumName')}
          placeholder="Enter album name"
        />
        
        <div className="flex items-center justify-end space-x-3 pt-4">
          <Button
            type="button"
            variant="ghost"
            onClick={onClose}
          >
            Cancel
          </Button>
          <Button
            type="submit"
            loading={isSubmitting}
          >
            Add Song
          </Button>
        </div>
      </form>
    </Modal>
  )
}

export default AddSongModal